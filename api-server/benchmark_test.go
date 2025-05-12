package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stellora/airline/api-server/api"
)

// BenchmarkSuite runs performance tests across different endpoints
// to establish a baseline for comparing microservices performance
func TestBenchmarkSuite(t *testing.T) {
	// Setup test environment
	ctx, handler := handlerTest(t)
	
	// Define benchmark cases for direct database access tests
	benchmarks := []struct {
		name   string
		action func(t *testing.T, ctx context.Context, handler *Handler) time.Duration
	}{
		{
			name: "ListAirlines",
			action: func(t *testing.T, ctx context.Context, handler *Handler) time.Duration {
				start := time.Now()
				v, err := handler.ListAirlines(ctx, api.ListAirlinesRequestObject{})
				if err != nil {
					t.Errorf("ListAirlines failed: %v", err)
					return 0
				}
				
				response, ok := v.(api.ListAirlines200JSONResponse)
				if !ok {
					t.Errorf("Expected ListAirlines200JSONResponse, got %T", v)
					return 0
				}
				t.Logf("Retrieved %d airlines", len(response))
				return time.Since(start)
			},
		},
		{
			name: "ListAirports",
			action: func(t *testing.T, ctx context.Context, handler *Handler) time.Duration {
				start := time.Now()
				v, err := handler.ListAirports(ctx, api.ListAirportsRequestObject{})
				if err != nil {
					t.Errorf("ListAirports failed: %v", err)
					return 0
				}
				
				response, ok := v.(api.ListAirports200JSONResponse)
				if !ok {
					t.Errorf("Expected ListAirports200JSONResponse, got %T", v)
					return 0
				}
				t.Logf("Retrieved %d airports", len(response))
				return time.Since(start)
			},
		},
		{
			name: "ListFlights",
			action: func(t *testing.T, ctx context.Context, handler *Handler) time.Duration {
				start := time.Now()
				v, err := handler.ListFlights(ctx, api.ListFlightsRequestObject{})
				if err != nil {
					t.Errorf("ListFlights failed: %v", err)
					return 0
				}
				
				response, ok := v.(api.ListFlights200JSONResponse)
				if !ok {
					t.Errorf("Expected ListFlights200JSONResponse, got %T", v)
					return 0
				}
				t.Logf("Retrieved %d flights", len(response))
				return time.Since(start)
			},
		},
	}

	// Insert test data
	insertAirlinesWithIATACodesT(t, handler, "XX", "YY")
	insertAirportsWithIATACodesT(t, handler, "AAA", "BBB", "CCC")
	insertAircraftT(t, handler, "XX", "REG1", "REG2")
	insertSchedulesT(t, handler, "XX1234 AAA-BBB")

	for _, bm := range benchmarks {
		t.Run(bm.name, func(t *testing.T) {
			// Execute test action and measure time
			duration := bm.action(t, ctx, handler)
			t.Logf("%s: response time %s", bm.name, duration)
		})
	}
}

// Actual Go benchmarks for precise performance measurements
func BenchmarkDataFetching(b *testing.B) {
	// Skip the benchmark in short mode
	if testing.Short() {
		b.Skip("skipping benchmark in short mode")
	}
	
	// Setup test environment
	ctx, handler := handlerTest(&testing.T{})
	
	// Configure handlers to benchmark
	benchmarks := []struct {
		name string
		fn   func(context.Context) error
	}{
		{
			name: "ListAirlines",
			fn: func(ctx context.Context) error {
				_, err := handler.ListAirlines(ctx, api.ListAirlinesRequestObject{})
				return err
			},
		},
		{
			name: "ListAirports",
			fn: func(ctx context.Context) error {
				_, err := handler.ListAirports(ctx, api.ListAirportsRequestObject{})
				return err
			},
		},
		{
			name: "ListFlights",
			fn: func(ctx context.Context) error {
				_, err := handler.ListFlights(ctx, api.ListFlightsRequestObject{})
				return err
			},
		},
		{
			name: "GetAirline",
			fn: func(ctx context.Context) error {
				_, err := handler.GetAirline(ctx, api.GetAirlineRequestObject{
					AirlineSpec: api.NewAirlineSpec(1, ""),
				})
				return err
			},
		},
		{
			name: "ListFlightsByAirline",
			fn: func(ctx context.Context) error {
				_, err := handler.ListFlightsByAirline(ctx, api.ListFlightsByAirlineRequestObject{
					AirlineSpec: api.NewAirlineSpec(1, ""),
				})
				return err
			},
		},
	}
	
	// Prepare test data
	t := &testing.T{}
	insertAirlinesWithIATACodesT(t, handler, "XX", "YY")
	insertAirportsWithIATACodesT(t, handler, "AAA", "BBB", "CCC")
	insertAircraftT(t, handler, "XX", "REG1", "REG2")
	insertSchedulesT(t, handler, "XX1234 AAA-BBB")
	if t.Failed() {
		b.Fatal("Failed to insert test data")
	}

	// Run benchmarks
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			// Reset timer for setup
			b.ResetTimer()
			
			// Run the benchmark
			for i := 0; i < b.N; i++ {
				if err := bm.fn(ctx); err != nil {
					b.Fatalf("%s failed: %v", bm.name, err)
				}
			}
		})
	}
}



// TestConcurrentDataFetching tests the ability to handle concurrent requests
func TestConcurrentDataFetching(t *testing.T) { 
	t.Skip("Skip concurrent test for now - needs database repair")
	// Skip test in short mode
	if testing.Short() {
		t.Skip("skipping concurrent test in short mode")
	}
	
	// Setup test environment
	ctx, handler := handlerTest(t)
	
	// Insert test data
	insertAirlinesWithIATACodesT(t, handler, "XX", "YY")
	insertAirportsWithIATACodesT(t, handler, "AAA", "BBB", "CCC")
	insertAircraftT(t, handler, "XX", "REG1", "REG2")
	insertSchedulesT(t, handler, "XX1234 AAA-BBB")
	
	// Define handlers to test concurrently
	handlers := []struct {
		name string
		fn   func(context.Context) error
	}{
		{
			name: "ListAirlines",
			fn: func(ctx context.Context) error {
				_, err := handler.ListAirlines(ctx, api.ListAirlinesRequestObject{})
				return err
			},
		},
		{
			name: "ListAirports",
			fn: func(ctx context.Context) error {
				_, err := handler.ListAirports(ctx, api.ListAirportsRequestObject{})
				return err
			},
		},
		{
			name: "ListFlights",
			fn: func(ctx context.Context) error {
				_, err := handler.ListFlights(ctx, api.ListFlightsRequestObject{})
				return err
			},
		},
		{
			name: "ListPassengers",
			fn: func(ctx context.Context) error {
				_, err := handler.ListPassengers(ctx, api.ListPassengersRequestObject{})
				return err
			},
		},
	}
	

	
	// Number of concurrent requests per handler
	concurrency := 10
	
	// Create channel for results
	results := make(chan string, len(handlers)*concurrency)
	
	// Launch concurrent requests
	for _, h := range handlers {
		for i := 0; i < concurrency; i++ {
			go func(handlerName string, fn func(context.Context) error) {
				err := fn(ctx)
				
				if err != nil {
					results <- fmt.Sprintf("%s: error %v", handlerName, err)
				} else {
					results <- fmt.Sprintf("%s: success", handlerName)
				}
			}(h.name, h.fn)
		}
	}
	
	// Collect results
	failures := 0
	for i := 0; i < len(handlers)*concurrency; i++ {
		result := <-results
		if result[len(result)-7:] != "success" {
			failures++
			t.Log(result)
		}
	}
	
	// Check if all requests were successful
	if failures > 0 {
		t.Errorf("%d concurrent requests failed", failures)
	}
}