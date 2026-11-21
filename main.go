package main

import "fmt"

func twoSum(nums []int, target int) []int {
	arr := make([]int, 0, 2)
	for i := 0; i < len(nums)-1; i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				arr = append(arr, i)
				arr = append(arr, j)
				return arr
			}
		}
	}
	return make([]int, 2)
}

func main() {
	arr := []int{2, 7, 9, 0}
	target := 9

	fmt.Println(twoSum(arr, target))
}
