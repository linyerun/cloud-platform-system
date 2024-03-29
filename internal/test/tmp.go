package test

import (
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
)

func main() {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
	// 查看镜像
	//images, err := client.ListImages(docker.ListImagesOptions{All: false})
	//if err != nil {
	//	panic(err)
	//}
	//for _, img := range images {
	//	fmt.Println("ID: ", img.ID)
	//	fmt.Println("RepoTags: ", img.RepoTags)
	//	fmt.Println("Created: ", img.Created)
	//	fmt.Println("Size: ", img.Size)
	//	fmt.Println("VirtualSize: ", img.VirtualSize)
	//	fmt.Println("ParentId: ", img.ParentID)
	//}

	// 拉取镜像
	err = client.PullImage(docker.PullImageOptions{Repository: "ubuntu", Tag: "latest"}, docker.AuthConfiguration{})
	if err != nil {
		panic(err)
	}

	// 查看镜像
	image, err := client.InspectImage("debian:latest")
	fmt.Println(image.ID, image.Size)

	//images, err := client.ListImages(docker.ListImagesOptions{Filter: "debian:latest", All: false})
	//if err != nil {
	//	panic(err)
	//}
	//for _, image := range images {
	//	fmt.Println(image.ID)
	//	fmt.Println(image.Size)
	//}

	// 删除镜像
	err = client.RemoveImage("sha256:5be4939875f7cdc771f11fb8c0737224da6d24e5e98023c511a5d1b4d1f94b04")
	if err != nil {
		panic(err)
	}
}
