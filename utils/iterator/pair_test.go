package iterutil

import (
	"fmt"
	"iter"
	"slices"
	"testing"
)

// func TestPullPair_BasicFunctionality(t *testing.T) {
// 	// Test with equal length slices
// 	iter1 := slices.All([]int{1, 2, 3})
// 	iter2 := slices.All([]string{"a", "b", "c"})
//
// 	transform := func(i int, s string) (string, error) {
// 		return fmt.Sprintf("%d-%s", i, s), nil
// 	}
//
// 	result := PullPair[int, string, string](iter1, iter2, transform)
// 	values := acc(result)
//
// 	expected := []any{"1-a", "2-b", "3-c"}
// 	if !slices.Equal(values, expected) {
// 		t.Errorf("expected %v, got %v", expected, values)
// 	}
// }
//
// func acc[T any](it iter.Seq2[T, error]) []any {
// 	var res []any
// 	for t, err := range it {
// 		if err != nil {
// 			panic(err)
// 		}
// 		res = append(res, t)
// 	}
// 	return res
// }

// Helper function to create an iterator from a slice with potential error
func sliceToIter[T any](slice []T, errorAt int) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for i, val := range slice {
			if i == errorAt {
				var zero T
				if !yield(zero, fmt.Errorf("error at index %d", i)) {
					return
				}
				return
			}
			if !yield(val, nil) {
				return
			}
		}
	}
}

// Helper function to collect results from an iterator
func collectResults[R any](seq iter.Seq2[R, error]) ([]R, error) {
	var results []R
	for val, err := range seq {
		if err != nil {
			return results, err
		}
		results = append(results, val)
	}
	return results, nil
}

func TestPullPair_BasicFunctionality(t *testing.T) {
	// Test with equal length slices
	iter1 := sliceToIter([]int{1, 2, 3}, -1)
	iter2 := sliceToIter([]string{"a", "b", "c"}, -1)

	transform := func(i int, s string) (string, error) {
		return fmt.Sprintf("%d-%s", i, s), nil
	}

	result := PullPair(iter1, iter2, transform)
	values, err := collectResults(result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"1-a", "2-b", "3-c"}
	if !slices.Equal(values, expected) {
		t.Errorf("expected %v, got %v", expected, values)
	}
}

func TestPullPair_EmptyIterators(t *testing.T) {
	iter1 := sliceToIter([]int{}, -1)
	iter2 := sliceToIter([]string{}, -1)

	transform := func(i int, s string) (string, error) {
		return fmt.Sprintf("%d-%s", i, s), nil
	}

	result := PullPair(iter1, iter2, transform)
	values, err := collectResults(result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(values) != 0 {
		t.Errorf("expected empty result, got %v", values)
	}
}

func TestPullPair_DifferentLengths(t *testing.T) {
	// First iterator longer
	iter1 := sliceToIter([]int{1, 2, 3, 4}, -1)
	iter2 := sliceToIter([]string{"a", "b", "c"}, -1)

	transform := func(i int, s string) (string, error) {
		return fmt.Sprintf("%d-%s", i, s), nil
	}

	result := PullPair(iter1, iter2, transform)
	_, err := collectResults(result)

	if err == nil {
		t.Fatal("expected error for different length iterators")
	}

	if err.Error() != "iterators have different lengths" {
		t.Errorf("expected 'iterators have different lengths', got %v", err)
	}

	// Second iterator longer
	iter1 = sliceToIter([]int{1, 2}, -1)
	iter2 = sliceToIter([]string{"a", "b", "c"}, -1)

	result = PullPair(iter1, iter2, transform)
	_, err = collectResults(result)

	if err == nil {
		t.Fatal("expected error for different length iterators")
	}

	if err.Error() != "iterators have different lengths" {
		t.Errorf("expected 'iterators have different lengths', got %v", err)
	}
}

func TestPullPair_ErrorInFirstIterator(t *testing.T) {
	iter1 := sliceToIter([]int{1, 2, 3}, 1) // Error at index 1
	iter2 := sliceToIter([]string{"a", "b", "c"}, -1)

	transform := func(i int, s string) (string, error) {
		return fmt.Sprintf("%d-%s", i, s), nil
	}

	result := PullPair(iter1, iter2, transform)
	values, err := collectResults(result)

	if err == nil {
		t.Fatal("expected error from first iterator")
	}

	if err.Error() != "error at index 1" {
		t.Errorf("expected 'error at index 1', got %v", err)
	}

	// Should have processed one value before error
	expected := []string{"1-a"}
	if !slices.Equal(values, expected) {
		t.Errorf("expected %v before error, got %v", expected, values)
	}
}

func TestPullPair_ErrorInSecondIterator(t *testing.T) {
	iter1 := sliceToIter([]int{1, 2, 3}, -1)
	iter2 := sliceToIter([]string{"a", "b", "c"}, 1) // Error at index 1

	transform := func(i int, s string) (string, error) {
		return fmt.Sprintf("%d-%s", i, s), nil
	}

	result := PullPair(iter1, iter2, transform)
	values, err := collectResults(result)

	if err == nil {
		t.Fatal("expected error from second iterator")
	}

	if err.Error() != "error at index 1" {
		t.Errorf("expected 'error at index 1', got %v", err)
	}

	// Should have processed one value before error
	expected := []string{"1-a"}
	if !slices.Equal(values, expected) {
		t.Errorf("expected %v before error, got %v", expected, values)
	}
}

func TestPullPair_ErrorInTransform(t *testing.T) {
	iter1 := sliceToIter([]int{1, 2, 3}, -1)
	iter2 := sliceToIter([]string{"a", "b", "c"}, -1)

	transform := func(i int, s string) (string, error) {
		if i == 2 {
			return "", fmt.Errorf("transform error at %d", i)
		}
		return fmt.Sprintf("%d-%s", i, s), nil
	}

	result := PullPair(iter1, iter2, transform)
	values, err := collectResults(result)

	if err == nil {
		t.Fatal("expected error from transform function")
	}

	if err.Error() != "transform error at 2" {
		t.Errorf("expected 'transform error at 2', got %v", err)
	}

	// Should have processed one value before error
	expected := []string{"1-a"}
	if !slices.Equal(values, expected) {
		t.Errorf("expected %v before error, got %v", expected, values)
	}
}

func TestPullPair_EarlyTermination(t *testing.T) {
	iter1 := sliceToIter([]int{1, 2, 3, 4, 5}, -1)
	iter2 := sliceToIter([]string{"a", "b", "c", "d", "e"}, -1)

	transform := func(i int, s string) (string, error) {
		return fmt.Sprintf("%d-%s", i, s), nil
	}

	result := PullPair(iter1, iter2, transform)

	var values []string
	count := 0
	for val, err := range result {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		values = append(values, val)
		count++
		if count == 2 {
			break // Early termination
		}
	}

	expected := []string{"1-a", "2-b"}
	if !slices.Equal(values, expected) {
		t.Errorf("expected %v, got %v", expected, values)
	}
}
