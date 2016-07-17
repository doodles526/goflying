package main

import (
	"fmt"
	"time"
	"github.com/westphae/goflying/mpu9250"
)

func main() {
	clock := time.NewTicker(100 * time.Millisecond)
	var (
		mpu      *mpu9250.MPU9250
		avg      *mpu9250.MPUData
		//cur      *mpu9250.MPUData
		err      error
	)

	for i:=0; i<10; i++ {
		mpu, err = mpu9250.NewMPU9250(250, 4, 50, true, false)
		if err != nil {
			fmt.Printf("Error initializing MPU9250, attempt %d of 10\n", i)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

	if err != nil {
		fmt.Println("Error: couldn't initialize MPU9250")
		return
	} else {
		fmt.Println("MPU9250 initialized successfully")
	}

	mpu.CCal<- 1
	fmt.Println("Awaiting Calibration Result")
	if err := <-mpu.CCalResult; err != nil {
		fmt.Println(err.Error())
		return
	} else {
		fmt.Println("Calibration succeeded")
	}

	for {
		<-clock.C

		avg = <-mpu.CAvg
		fmt.Printf("\nTime:   %6.1f ms\n", float64(avg.DT.Nanoseconds())/1000000)
		fmt.Printf("Number of Observations: %d\n", avg.N)
		fmt.Printf("Avg Gyro:   % +8.2f % +8.2f % +8.2f\n", avg.G1, avg.G2, avg.G3)
		fmt.Printf("Avg Accel:  % +8.2f % +8.2f % +8.2f\n", avg.A1, avg.A2, avg.A3)

		if !mpu.MagEnabled() {
			fmt.Println("Magnetometer disabled")
		} else if avg.MagError != nil {
			fmt.Println(avg.MagError.Error())
		} else {
			fmt.Printf("Mag:        % +8.0f % +8.0f % +8.0f\n", avg.M1, avg.M2, avg.M3)
		}

		/*
		cur = <-mpu.C
		fmt.Printf("Cur Gyro:   % +8.2f % +8.2f % +8.2f\n", cur.G1, cur.G2, cur.G3)
		fmt.Printf("Cur Accel:  % +8.2f % +8.2f % +8.2f\n", cur.A1, cur.A2, cur.A3)

		fmt.Printf("Length of buffered channel: %d\n", len(mpu.CBuf))
		for i := 0; i<10; i++ {
			d := <-mpu.CBuf
			fmt.Printf("%6.3f ", float64(d.T.Nanosecond())/1000000000)
		}
		fmt.Println()
		*/
	}
}
