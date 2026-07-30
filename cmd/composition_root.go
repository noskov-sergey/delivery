package cmd

import (
	"database/sql"
	"delivery/internal/adapters/out/grpc/geo"
	"delivery/internal/adapters/out/postgres"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"delivery/internal/jobs"
	"log"
	"sync"

	"github.com/robfig/cron/v3"
)

type CompositionRoot struct {
	configs Config

	db  *sql.DB
	uow ports.UnitOfWork

	geoClient ports.GeoClient
	onceGeo   sync.Once

	closers []Closer
}

func NewCompositionRoot(configs Config, db *sql.DB) *CompositionRoot {
	return &CompositionRoot{
		configs: configs, db: db,
	}
}

///////////////////////////////////////////////////////////
//////////////////// LIFECYCLE ////////////////////////////
///////////////////////////////////////////////////////////

func (cr *CompositionRoot) RegisterCloser(c Closer) {
	cr.closers = append(cr.closers, c)
}

func (cr *CompositionRoot) CloseAll() {
	for _, c := range cr.closers {
		if err := c.Close(); err != nil {
			log.Printf("close error: %v", err)
		}
	}
}

func (cr *CompositionRoot) NewOrderDispatcherService() services.OrderDispatchService {
	orderDispatcher := services.NewOrderDispatchService()
	return orderDispatcher
}

func (cr *CompositionRoot) NewUnitOfWork() ports.UnitOfWork {
	if cr.uow != nil {
		return cr.uow
	}

	uow, err := postgres.NewUnitOfWork(cr.db)
	if err != nil {
		return nil
	}

	cr.uow = uow

	return uow
}

func (cr *CompositionRoot) NewCreateOrderCommandHandler() commands.CreateOrderCommandHandler {
	h, err := commands.NewCreateOrderCommandHandler(cr.NewUnitOfWork(), cr.NewGeoClient())
	if err != nil {
		log.Fatalf("ERROR: cannot create CreateOrderCommandHandler: %v", err)
	}
	return h
}

func (cr *CompositionRoot) NewCreateCourierCommandHandler() commands.CreateCourierCommandHandler {
	h, err := commands.NewCreateCourierCommandHandler(cr.NewUnitOfWork())
	if err != nil {
		log.Fatalf("ERROR: cannot create CreateCourierCommandHandler: %v", err)
	}
	return h
}

func (cr *CompositionRoot) NewGetAllCouriersQueryHandler() queries.GetAllQueriersQueryHandler {
	h, err := queries.NewGetAllQueriersQueryHandler(cr.NewUnitOfWork())
	if err != nil {
		log.Fatalf("ERROR: cannot create GetAllCouriersQueryHandler: %v", err)
	}
	return h
}

func (cr *CompositionRoot) NewGetIncompletedOrdersQueryHandler() queries.GetAllUncompletedOrdersQueryHandler {
	h, err := queries.NewGetAllUncompletedOrderQueryHandler(cr.NewUnitOfWork())
	if err != nil {
		log.Fatalf("ERROR: cannot create GetIncompletedOrdersQueryHandler: %v", err)
	}
	return h
}

func (cr *CompositionRoot) NewAssignOrdersCommandHandler() commands.AssignOrdersCommandHandler {
	h, err := commands.NewAssignOrdersCommandHandler(cr.NewUnitOfWork(), cr.NewOrderDispatcherService())
	if err != nil {
		log.Fatalf("ERROR: cannot create AssignOrdersCommandHandler: %v", err)
	}
	return h
}

func (cr *CompositionRoot) NewMoveCouriersCommandHandler() commands.MoveCouriersCommandHandler {
	h, err := commands.NewMoveCouriersCommandHandler(cr.NewUnitOfWork())
	if err != nil {
		log.Fatalf("ERROR: cannot create AssignMoveCouriersCommandHandler: %v", err)
	}
	return h
}

func (cr *CompositionRoot) NewAssignOrdersJob() cron.Job {
	job, err := jobs.NewAssignOrdersJob(cr.NewAssignOrdersCommandHandler())
	if err != nil {
		log.Fatalf("cannot create AssignOrdersJob: %v", err)
	}
	return job
}
func (cr *CompositionRoot) NewMoveCouriersJob() cron.Job {
	job, err := jobs.NewMoveCouriersJob(cr.NewMoveCouriersCommandHandler())
	if err != nil {
		log.Fatalf("cannot create MoveCouriersJob: %v", err)
	}
	return job
}

func (cr *CompositionRoot) NewGeoClient() ports.GeoClient {
	cr.onceGeo.Do(func() {
		client, err := geo.NewClient(cr.configs.GeoServiceGrpcHost)
		if err != nil {
			log.Fatalf("ERROR: create GeoClient: %v", err)
		}
		cr.RegisterCloser(client)
		cr.geoClient = client
	})
	return cr.geoClient
}
