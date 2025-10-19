package main

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	metricv1 "github.com/naivary/omp/api/metric/v1"
)

func CreateMetricDefinition(pg *pgxpool.Pool) *Endpoint {
	hl := HandlerFuncErr(func(w http.ResponseWriter, r *http.Request) error {
		data, err := decode[metricv1.CreateMetricRequest](r)
		if err != nil {
			return err
		}
		fmt.Println(data)
		return encode[any](w, r, http.StatusCreated, nil)
	})

	return &Endpoint{
		Handler: hl,
	}
}
