package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func init() {
	fmt.Println("OSS Go SDK Version: ", oss.Version)
}

func handleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}
func ossTest() {
	// 创建OSSClient实例。
	client, err := oss.New("http://oss-cn-beijing.aliyuncs.com", "LTAI4Fjt9mcQoUkx8nDmmJSV", "mY03PEw5lUOGCM6upTEtD11yHzi7gF")
	if err != nil {
		handleError(err)
	}
	// 获取存储空间。
	bucketName := "jw397"
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		handleError(err)
	}

	// 获取存储空间的信息，包括地域（Region或Location）、创建日期（CreationDate）、访问权限（ACL）、拥有者（Owner）、存储类型（StorageClass）、容灾类型（RedundancyType）等。
	// res, err := client.GetBucketInfo(bucketName)
	// if err != nil {
	// 	handleError(err)
	// }
	// fmt.Println("BucketInfo: ", res.BucketInfo)
	// fmt.Println("BucketInfo.Location: ", res.BucketInfo.Location)
	// fmt.Println("BucketInfo.CreationDate: ", res.BucketInfo.CreationDate)
	// fmt.Println("BucketInfo.ACL: ", res.BucketInfo.ACL)
	// fmt.Println("BucketInfo.Owner: ", res.BucketInfo.Owner)
	// fmt.Println("BucketInfo.StorageClass: ", res.BucketInfo.StorageClass)
	// fmt.Println("BucketInfo.RedundancyType: ", res.BucketInfo.RedundancyType)
	// fmt.Println("BucketInfo.ExtranetEndpoint: ", res.BucketInfo.ExtranetEndpoint)
	// fmt.Println("BucketInfo.IntranetEndpoint: ", res.BucketInfo.IntranetEndpoint)

	// 列举文件。
	// marker := ""
	// for {
	// 	lsRes, err := bucket.ListObjects(oss.Marker(marker))
	// 	if err != nil {
	// 		handleError(err)
	// 	}
	// 	// 打印列举文件，默认情况下一次返回100条记录。
	// 	for _, object := range lsRes.Objects {
	// 		fmt.Println("Bucket: ", object.Key)
	// 	}
	// 	if lsRes.IsTruncated {
	// 		marker = lsRes.NextMarker
	// 	} else {
	// 		break
	// 	}
	// }

	// 列举指定前缀的存储空间。
	// lsRes, err := client.ListBuckets(oss.Prefix("download/"))
	// if err != nil {
	// 	handleError(err)
	// }
	// // 打印存储空间列表。
	// fmt.Println("Buckets with prefix: ", lsRes.Buckets)
	// for _, bucket := range lsRes.Buckets {
	// 	fmt.Println("Bucket with prefix: ", bucket.Name)
	// }

	// 读取本地文件。
	fd, err := http.Get("https://www.7-zip.org/a/7z1900.exe")
	if err != nil {
		handleError(err)
	}
	// fmt.Println(fd)
	defer fd.Body.Close()

	// 上传文件流。
	err = bucket.PutObject("jw/test", fd.Body)
	if err != nil {
		handleError(err)
	}
}

// func testGetFileType() {
// 	f, err := http.Head("https://www.7-zip.org/a/7z1900.exe")
// 	if err != nil {
// 		fmt.Println("open error: %v", err)
// 	}
// 	// defer f.Body.Close()
// 	// fSrc, err := ioutil.ReadAll(f)
// 	// t.Log(GetFileType(fSrc[:10]))
// 	fmt.Println(f)
// }

func main() {
	// testGetFileType()

	ossTest()

	// fd, err := http.Get("https://www.7-zip.org/a/7z1900.exe")
	// if err != nil {
	// 	handleError(err)
	// }
	// fmt.Println(fd.Body)
}
