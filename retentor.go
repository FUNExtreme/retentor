package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	minio "github.com/minio/minio-go"
)

// DigitalOcean specifies a DigitalOcean space.
type digitalOcean struct {
	// Available regions:
	// * nyc3
	// * ams3
	// * sgp1
	Region string
	// Name of the space
	Space string
	// DigitalOcean Spaces access key
	AccessKey string
	// DigitalOcean Spaces access secret
	ClientSecret string
}

type spaceObject struct {
	path       string
	pathChunks []string
	filename   string
	dateUtc    string
	timeUtc    string
}

type pathDescriptionToIndex struct {
	customer int
	dateUtc  int
	timeUtc  int
}

func main() {
	// DigitalOcean credentials
	doSpace := os.Getenv("DO_SPACE")
	doSpaceRegion := os.Getenv("DO_SPACE_REGION")
	doSpaceAccessKey := os.Getenv("DO_SPACE_ACCESS_KEY")
	doSpaceClientSecret := os.Getenv("DO_SPACE_CLIENT_SECRET")

	do := &digitalOcean{
		Region:       doSpaceRegion,
		Space:        doSpace,
		AccessKey:    doSpaceAccessKey,
		ClientSecret: doSpaceClientSecret,
	}

	pathDescriptor := &pathDescriptionToIndex{
		customer: 0,
		dateUtc:  1,
		timeUtc:  -1,
	}

	ssl := true
	client, err := minio.New(do.Region+".digitaloceanspaces.com", do.AccessKey, do.ClientSecret, ssl)
	if err != nil {
		log.Fatal(err)
	}

	// Create a done channel to control 'ListObjects' go routine.
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	var spaceObjects []spaceObject

	// List all objects from a bucket-name with a matching prefix.
	for object := range client.ListObjectsV2(do.Space, "", true, doneCh) {
		if object.Err != nil {
			fmt.Println(object.Err)
			return
		}

		pathSplit := strings.Split(object.Key, "/")
		spaceObj := &spaceObject{
			path:       object.Key,
			pathChunks: pathSplit,
			filename:   pathSplit[len(pathSplit)-1],
		}
		if pathDescriptor.dateUtc != -1 {
			spaceObj.dateUtc = pathSplit[pathDescriptor.dateUtc]
		}
		if pathDescriptor.timeUtc != -1 {
			spaceObj.timeUtc = pathSplit[pathDescriptor.timeUtc]
		}
		spaceObjects = append(spaceObjects, *spaceObj)
	}

	spaceObjectsByDate := make(map[string][]spaceObject)
	for _, obj := range spaceObjects {
		spaceObjectsByDate[obj.dateUtc] = append(spaceObjectsByDate[obj.dateUtc], obj)
	}

	currentTime := time.Now().UTC()
	for _, retentionPolicy := range RetentionPolicies {
		fmt.Println("------------------------")
		fmt.Println("policy:", retentionPolicy.PolicyType)
		fmt.Println("---")
		startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).AddDate(0, 0, -retentionPolicy.StartsAfterDays)
		endTime := startTime.AddDate(0, 0, -retentionPolicy.ValidForDays)
		fmt.Println("start (inclusive):", startTime.Format("20060102"))
		fmt.Println("end (exclusive):", endTime.Format("20060102"))
		fmt.Println("---")

		policyObjects := []spaceObject{}
		for date, objects := range spaceObjectsByDate {
			objDate, _ := time.Parse("20060102", date)
			if (startTime.Equal(objDate) || startTime.After(objDate)) && endTime.Before(objDate) {
				for _, obj := range objects {
					policyObjects = append(policyObjects, obj)
					fmt.Println("Object within policy date range: " + obj.path)
				}
			}
		}

		//markedForDeletionObjects := []spaceObject{}
		switch retentionPolicy.PolicyType {
		//case RetentionPolicyHourly:
		//case RetentionPolicySixHourly:
		//case RetentionPolicyDaily:
		//case RetentionPolicyWeekly:
		//case RetentionPolicyMonthly:
		default:
			fmt.Println("No code for this retention policy")
		}
	}

	return
}
