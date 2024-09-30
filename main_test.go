package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"
	"time"

	"github.com/MukeshGKastala/marketing/api"
	sqlc "github.com/MukeshGKastala/marketing/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestSequentialCreatePromotionRequests(t *testing.T) {
	t.Parallel()

	db := createTestDB(t, "8.0")
	defer db.close()

	require.NoError(t, runMigration(db.db))

	svc := api.NewServer(sqlc.NewStore(db.db))

	for i := 0; i < 5; i++ {
		resp, err := svc.CreatePromotion(context.Background(), api.CreatePromotionRequestObject{
			Body: &api.CreatePromotionJSONRequestBody{
				PromotionCode: "BOGO50",
				StartDate:     time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
				EndDate:       time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC),
			},
		})

		require.NoError(t, err)
		w := httptest.NewRecorder()
		require.NoError(t, resp.VisitCreatePromotionResponse(w))
		var want int
		if i == 0 {
			want = http.StatusCreated
		} else {
			want = http.StatusBadRequest
		}
		require.Equal(t, want, w.Code)
	}
}

func TestConcurrentCreatePromotionRequests(t *testing.T) {
	t.Parallel()

	db := createTestDB(t, "8.0")
	defer db.close()

	require.NoError(t, runMigration(db.db))

	svc := api.NewServer(sqlc.NewStore(db.db))

	g, ctx := errgroup.WithContext(context.Background())

	for i := 0; i < 5; i++ {
		g.Go(func() error {
			_, err := svc.CreatePromotion(ctx, api.CreatePromotionRequestObject{
				Body: &api.CreatePromotionJSONRequestBody{
					PromotionCode: "BOGO50",
					StartDate:     time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
					EndDate:       time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC),
				},
			})
			return err
		})
	}

	require.NoError(t, g.Wait())

	obj, err := svc.ListPromotions(context.Background(), api.ListPromotionsRequestObject{})
	require.NoError(t, err)
	w := httptest.NewRecorder()
	require.NoError(t, obj.VisitListPromotionsResponse(w))
	res := w.Result()
	b, err := io.ReadAll(res.Body)
	var promotions api.ListPromotions200JSONResponse
	require.NoError(t, json.Unmarshal(b, &promotions))
	require.Len(t, promotions, 1)
}

// testDB is a wrapper around *sql.DB for MySQL connections
// designed for tests to provide a clean database for each testcase.
// Callers should cleanup with close() when finished.
type testDB struct {
	db        *sql.DB
	container *dockertest.Resource
}

// close releases the resources used by the Pool and
// removes the container and linked volumes from docker.
// Failure to do so will result in dangling containers.
func (t *testDB) close() error {
	t.db.Close()
	return t.container.Close()
}

// createTestDB returns a testDB which can be used in tests as a clean
// MySQL database. Callers should call Close() on the returned *testDB.
func createTestDB(t *testing.T, version string) *testDB {
	if !dockerEnabled() {
		t.Skip("docker not enabled")
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		t.Fatalf("could not connect to Docker: %s", err)
	}

	resource, err := pool.Run("mysql", version, []string{"MYSQL_ROOT_PASSWORD=secret", "MYSQL_DATABASE=marketing"})
	if err != nil {
		t.Fatalf("could not start mysql resource: %s", err)
	}

	var db *sql.DB
	dsn := fmt.Sprintf("root:secret@tcp(127.0.0.1:%s)/marketing?multiStatements=true&parseTime=true", resource.GetPort("3306/tcp"))

	// Exponential backoff for up to 1 minute
	if err = pool.Retry(func() error {
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		resource.Close()
		t.Fatalf("could not connect to mysql container: %s", err)
	}

	return &testDB{db, resource}
}

// dockerEnabled returns true if Docker is available when called.
func dockerEnabled() bool {
	bin, err := exec.LookPath("docker")
	return bin != "" && err == nil // 'docker' was found on PATH
}
