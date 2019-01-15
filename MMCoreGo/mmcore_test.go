package mmcore_test

import (
	"fmt"
	"log"

	mmcore "github.com/Andeling/MMCoreAPI/MMCoreGo"
)

func ExampleSession() {
	mmc := mmcore.NewSession()
	defer mmc.Close()

	fmt.Printf("Version Info: %s\n", mmc.VersionInfo())
	fmt.Printf("API Version Info: %s\n", mmc.APIVersionInfo())
	// Output:
	// Version Info: MMCore version 8.6.0
	// API Version Info: Device API version 68, Module API version 10
}

func ExampleSession_SnapImage() {
	mmc := mmcore.NewSession()
	defer mmc.Close()

	// Set the search path for device adapters
	// MMCore will use "mmgr_dal_DemoCamera.dll" when we load a device with DemoCamera module.
	mmc.SetDeviceAdapterSearchPaths([]string{"C:\\Program Files\\Micro-Manager-1.4"})

	// MMCore refers a device by the label name.
	cameraLabel := "Camera"

	// Load device "DCam" with "DemoCamera" module, and assign the label name.
	err := mmc.LoadDevice(cameraLabel, "DemoCamera", "DCam")
	if err != nil {
		log.Fatal(err)
	}

	//
	err = mmc.InitializeAllDevices()
	if err != nil {
		log.Fatal(err)
	}

	// Set the camera as default camera device in the session.
	// SnapImage() and StartContinuousSequenceAcquisition() can only use the default camera.
	// The MMCore C++ API allows StartSequenceAcquisition to use non-default camera,
	// but that is not implemented in the MMCoreC and MMCoreGo API.
	// To access multiple cameras, just create a session for each of the cameras.
	err = mmc.SetCameraDevice(cameraLabel)
	if err != nil {
		log.Fatal(err)
	}

	// Get the parameters which we will need to interprete the raw image data.
	// These will not change without a SetProperty.
	width := mmc.ImageWidth()
	height := mmc.ImageHeight()
	bytesPerPixel := mmc.BytesPerPixel()
	bitDepth := mmc.ImageBitDepth()
	fmt.Printf("width=%d, height=%d, bytesPerPixel=%d, bitDepth=%d\n", width, height, bytesPerPixel, bitDepth)

	// Set exposure time to 100 ms
	err = mmc.SetExposureTime(100)
	if err != nil {
		log.Fatal(err)
	}

	// Get the acutal exposure time
	// When controlling an actual hardware,
	// this number will not be exactly the same as the requested exposure time,
	// and it will be some multiply of the internal clock of the camera.
	// It is good to get and save the actual exposure time when acquiring an image.
	exposure, err := mmc.ExposureTime()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Exposure time: %.6f ms\n", exposure)

	// Start the exposure and wait for the exposure to finish
	err = mmc.SnapImage()
	if err != nil {
		log.Fatal(err)
	}

	// Wait for the read-out and data transfering and get the image data.
	// This may take a while in the case of an actual camera.
	// You will not be able to start another exposure with SnapImage in the meantime.
	//
	buf, err := mmc.GetImage()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("len(buf)=%d\n", len(buf))

	// Output:
	// width=512, height=512, bytesPerPixel=1, bitDepth=8
	// Exposure time: 100.000000 ms
	// len(buf)=262144
}

func ExampleSession_StartContinuousSequenceAcquisition() {
	mmc := mmcore.NewSession()
	defer mmc.Close()

	// Set the search path for device adapters
	// MMCore will use "mmgr_dal_DemoCamera.dll" when we load a device with DemoCamera module.
	mmc.SetDeviceAdapterSearchPaths([]string{"C:\\Program Files\\Micro-Manager-1.4"})

	// MMCore refers a device by the label name.
	cameraLabel := "Camera"

	// Load device "DCam" with "DemoCamera" module, and assign the label name.
	err := mmc.LoadDevice(cameraLabel, "DemoCamera", "DCam")
	if err != nil {
		log.Fatal(err)
	}

	//
	err = mmc.InitializeAllDevices()
	if err != nil {
		log.Fatal(err)
	}

	// Set the camera as default camera device in the session.
	// SnapImage() and StartContinuousSequenceAcquisition() can only use the default camera.
	// The MMCore C++ API allows StartSequenceAcquisition to use non-default camera,
	// but that is not implemented in the MMCoreC and MMCoreGo API.
	// To access multiple cameras, just create a session for each of the cameras.
	err = mmc.SetCameraDevice(cameraLabel)
	if err != nil {
		log.Fatal(err)
	}

	// Get the parameters which we will need to interprete the raw image data.
	// These will not change without a SetProperty.
	width := mmc.ImageWidth()
	height := mmc.ImageHeight()
	bytesPerPixel := mmc.BytesPerPixel()
	bitDepth := mmc.ImageBitDepth()
	fmt.Printf("width=%d, height=%d, bytesPerPixel=%d, bitDepth=%d\n", width, height, bytesPerPixel, bitDepth)

	// Set exposure time to 100 ms
	err = mmc.SetExposureTime(100)
	if err != nil {
		log.Fatal(err)
	}

	// Get the acutal exposure time
	// When controlling an actual hardware,
	// this number will not be exactly the same as the requested exposure time,
	// and it will be some multiply of the internal clock of the camera.
	// It is good to get and save the actual exposure time when acquiring an image.
	exposure, err := mmc.ExposureTime()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Exposure time: %.6f ms\n", exposure)

	err = mmc.StartContinuousSequenceAcquisition(0)
	if err != nil {
		log.Fatal(err)
	}

	n_images := 0
	for {
		if mmc.GetRemainingImageCount() > 0 {
			_, err := mmc.GetLastImage()
			if err != nil {
				mmc.StopSequenceAcquisition()
				log.Fatal(err)
			}
			n_images++
			if n_images == 10 {
				break
			}
		}
	}

	err = mmc.StopSequenceAcquisition()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Finished acquiring %d images with ContinuousSequenceAcquisition.\n", n_images)

	// Output:
	// width=512, height=512, bytesPerPixel=1, bitDepth=8
	// Exposure time: 100.000000 ms
	// Finished acquiring 10 images with ContinuousSequenceAcquisition.
}