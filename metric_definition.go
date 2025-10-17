package main

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"

	metricv1 "github.com/naivary/omp/api/metric/v1"
)

func CreateMetricDefinition(pg *pgx.Conn) *Endpoint {
	hl := HandlerFuncErr(func(w http.ResponseWriter, r *http.Request) error {
		data, err := decode[metricv1.CreateMetricRequest](r)
		if err != nil {
			return err
		}
		def := metricv1.Definition{}
		fmt.Println(data)
		return encode(w, r, http.StatusCreated, def)
	})

	return &Endpoint{
		Handler: hl,
	}
}
