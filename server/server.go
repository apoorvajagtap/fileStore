package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/apoorvajagtap/fileStore/dbase"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type File struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	FileSize   int64              `json:"filesize" bson:"length"`
	ChunkSize  int64              `json:"chunksize" bson:"chunkSize"`
	UploadDate time.Time          `json:"uploaddate" bson:"uploadDate"`
	Name       string             `json:"name" bson:"filename"`
}

type Chunk struct {
	ChunkId  primitive.ObjectID `json:"cid" bson:"_id"`
	FileId   primitive.ObjectID `json:"fid" bson:"files_id"`
	Sequence int                `json:"seq" bson:"n"`
	Content  bson.RawValue      `json:"content" bson:"data"`
}

var fileCollection *mongo.Collection = dbase.GetCollection(dbase.DB, "fileStore_collection", "fileStore_db")

func fileExists(fileHeader *multipart.FileHeader) bool {

	var results bson.M

	fsFiles := fileCollection.Database().Collection("fs.files")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	_ = fsFiles.FindOne(ctx, bson.D{{Key: "filename", Value: fileHeader.Filename}}).Decode(&results)
	return results != nil
}

func returnFileList() []File {
	var results []File
	findOptions := options.Find()

	fsFiles := fileCollection.Database().Collection("fs.files")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	cursor, err := fsFiles.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cursor.Next(ctx) {
		var res File
		err := cursor.Decode(&res)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, res)
	}

	cursor.Close(context.TODO())
	return results
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
		return
	}

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["file"]

	// create bucket
	bucket, err := gridfs.NewBucket(
		fileCollection.Database(),
	)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	for _, fileHeader := range files {

		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// check if file with same name already exists
		if fileExists(fileHeader) {
			http.Error(w, fmt.Sprintf(">>> The file '%s' already exists!\n", fileHeader.Filename), http.StatusInternalServerError)
			log.Printf("The file '%s' already exists!\n", fileHeader.Filename)
			continue
		}

		// saving the file data in buffer
		buff := bytes.NewBuffer(nil)
		if _, err := io.Copy(buff, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		uploadStream, err := bucket.OpenUploadStream(
			fileHeader.Filename,
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer uploadStream.Close()

		fileSize, err := uploadStream.Write([]byte(buff.String()))
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		log.Printf("Write file to DB was successful. File size: %d\n", fileSize)
		fmt.Fprintf(w, ">>> File %s uploaded successfully", uploadStream.FileID)
	}
}

// list all the fileNames
func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
		return
	}

	fileList := returnFileList()

	for _, f := range fileList {
		w.Write([]byte(fmt.Sprintf("%s\n", f.Name)))
		fmt.Println(f)
	}
}

// Deletes single file at a time.
func deleteFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
		return
	}

	// find the ID of the file with given name
	var fileResult File

	fsFiles := fileCollection.Database().Collection("fs.files")
	fsChunks := fileCollection.Database().Collection("fs.chunks")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err := fsFiles.FindOne(ctx, bson.D{{Key: "filename", Value: r.Header["Name"][0]}}).Decode(&fileResult)
	if err != nil {
		log.Panic(err)
	}

	// deleting from fs.chunks (content)
	_, err = fsChunks.DeleteOne(ctx, bson.D{{Key: "files_id", Value: fileResult.Id}}, nil)
	if err != nil {
		log.Panic(err)
		return
	}

	// deleting from fs.files (metadata)
	_, err = fsFiles.DeleteOne(ctx, bson.D{{Key: "_id", Value: fileResult.Id}}, nil)
	if err != nil {
		log.Panic(err)
		return
	}

	w.Write([]byte(fmt.Sprintf("'%s' has been deleted successfully! ", fileResult.Name)))
	log.Printf("File '%s' with id: '%s' has been deleted successfully! \n", fileResult.Name, fileResult.Id)
}

func modifyFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	// saving the file data in buffer
	buff := bytes.NewBuffer(nil)
	if _, err := io.Copy(buff, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var fileResult File

	fsFiles := fileCollection.Database().Collection("fs.files")
	fsChunks := fileCollection.Database().Collection("fs.chunks")

	err = fsFiles.FindOne(context.TODO(), bson.D{{Key: "filename", Value: fileHeader.Filename}}).Decode(&fileResult)
	if err != nil {
		log.Panic(err)
	}

	filter := bson.D{{Key: "files_id", Value: fileResult.Id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "content", Value: buff.String()}}}}

	opts := options.Update().SetUpsert(true)

	result, err := fsChunks.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Panic(err)
	}

	if result.ModifiedCount == 0 {
		w.Write([]byte("No changes in the content observed!"))
	} else {
		w.Write([]byte(fmt.Sprintf("Modified the content of %s", fileHeader.Filename)))
	}
	fmt.Printf("Number of documents updated: %v\n", result.ModifiedCount)
	fmt.Printf("Number of documents upserted: %v\n", result.UpsertedCount)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", uploadFileHandler)
	mux.HandleFunc("/get", listFilesHandler)
	mux.HandleFunc("/delete", deleteFileHandler)
	mux.HandleFunc("/update", modifyFileHandler)

	if err := http.ListenAndServe(":4500", mux); err != nil {
		log.Fatal(err)
	}
}
