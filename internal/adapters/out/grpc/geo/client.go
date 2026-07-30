package geo

import (
	"context"
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/ports"
	"delivery/internal/generated/clients/geopb"
	"errors"
	"time"

	"github.com/labstack/gommon/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ ports.GeoClient = &client{}

type client struct {
	conn        *grpc.ClientConn
	pbGeoClient geopb.GeoClient
	timeout     time.Duration
}

func NewClient(host string) (ports.GeoClient, error) {
	if host == "" {
		return nil, errors.New("host required")
	}

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err.Error())
	}

	pbGeoClient := geopb.NewGeoClient(conn)

	return &client{
		conn:        conn,
		pbGeoClient: pbGeoClient,
		timeout:     time.Second * 5,
	}, nil
}

func (c *client) GetLocation(ctx context.Context, street string) (kernel.Location, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.pbGeoClient.GetGeolocation(ctx, &geopb.GetGeolocationRequest{
		Street: street,
	})
	if err != nil {
		return kernel.Location{}, err
	}

	return kernel.NewLocation(uint8(resp.Location.GetX()), uint8(resp.Location.GetY()))
}

func (c *client) Close() error {
	return c.conn.Close()
}
