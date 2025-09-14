package handler

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"referral-service/controller"
	"referral-service/domain"

	pb "referral-service/proto/referral"
)

type Handlers struct {
	pb.UnimplementedReferralServiceServer

	log         *zap.Logger
	referralCon controller.ReferralController
	programCon  controller.ProgramController
	memberCon   controller.MemberController
	health      *health.Server
}

// Params defines constructor requirements.
type Params struct {
	fx.In

	Log         *zap.Logger
	Lc          fx.Lifecycle
	Cfg         config.Provider
	ReferralCon controller.ReferralController
	ProgramCon  controller.ProgramController
	MemberCon   controller.MemberController
}

// New is the handler constructor.
func New(p Params) (*Handlers, error) {
	h := &Handlers{
		log:         p.Log,
		referralCon: p.ReferralCon,
		programCon:  p.ProgramCon,
		memberCon:   p.MemberCon,
	}
	ln, err := net.Listen(
		"tcp",
		":5000",
	)
	if err != nil {
		return nil, fmt.Errorf("grpc net listen %w", err)
	}

	// Create grpc server.
	grpcServer := grpc.NewServer()

	// Add reflection to service stack.
	reflection.Register(grpcServer)

	// Add healthcheck to service stack.
	healthCheck := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthCheck)
	h.health = healthCheck

	// Add sample proto service to service stack.
	pb.RegisterReferralServiceServer(grpcServer, h)

	// gRPC client connection for HTTP proxy.
	conn, err := grpc.DialContext(
		context.Background(),
		":5000",
		// grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc dial context %w", err)
	}

	// Create proxy.
	gwmux := runtime.NewServeMux()
	// Register proxy handlers. Routes http calls to gRPC.
	err = pb.RegisterReferralServiceHandler(
		context.Background(),
		gwmux,
		conn,
	)
	if err != nil {
		return nil, fmt.Errorf("register proxy handler %w", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	p.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Start gRPC server.
			h.log.Info("Serving gRPC",
				zap.String("address", p.Cfg.Get("server.address").String()),
			)
			go func() {
				if err := grpcServer.Serve(ln); err != nil {
					h.log.Error("grpc serve", zap.Error(err))
					return
				}
			}()

			// Start proxy server.
			h.log.Info("Starting http proxy", zap.String("address", gwServer.Addr))
			go func() {
				if err := gwServer.ListenAndServe(); err != nil {
					h.log.Error("proxy listen&serve", zap.Error(err))
					return
				}
			}()

			// Set initial health status.
			h.health.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			h.log.Info("shutting down")
			grpcServer.GracefulStop()
			gwServer.Shutdown(ctx)
			return nil
		},
	})

	return h, nil
}

// -------------------------------------------------------------
// Program API handlers
// -------------------------------------------------------------

func (h *Handlers) GetPrograms(
	ctx context.Context,
	req *pb.GetProgramsRequest,
) (*pb.GetProgramsResponse, error) {
	var page = 1
	if req.Page != nil {
		page = int(*req.Page)
	}
	var size = 100
	if req.Size != nil {
		size = int(*req.Size)
	}

	programs, err := h.programCon.GetPrograms(ctx, page, size)
	if err != nil {
		return &pb.GetProgramsResponse{}, err
	}

	protoPrograms := make([]*pb.Program, 0, len(programs))
	for _, p := range programs {
		protoPrograms = append(protoPrograms, ToProtoProgram(p))
	}

	return &pb.GetProgramsResponse{
		Programs: protoPrograms,
	}, nil
}

func (h *Handlers) GetProgram(
	ctx context.Context,
	req *pb.GetProgramRequest,
) (*pb.GetProgramResponse, error) {
	program, err := h.programCon.GetProgram(ctx, req.Id)
	if err != nil {
		return &pb.GetProgramResponse{}, err
	}

	return &pb.GetProgramResponse{
		Program: ToProtoProgram(*program),
	}, nil
}

func (h *Handlers) AddProgram(
	ctx context.Context,
	req *pb.AddProgramRequest,
) (*pb.AddProgramResponse, error) {
	programId, err := h.programCon.AddProgram(ctx, req.Name, req.Title, req.Active)

	if err != nil {
		return &pb.AddProgramResponse{}, err
	}

	return &pb.AddProgramResponse{
		Id: programId,
	}, nil
}

func (h *Handlers) UpdateProgram(
	ctx context.Context,
	req *pb.UpdateProgramRequest,
) (*pb.UpdagteProgramResponse, error) {
	program, err := h.programCon.UpdateProgram(ctx, req.Id, req.Name, req.Title, req.Active)

	if err != nil {
		return &pb.UpdagteProgramResponse{}, err
	}

	return &pb.UpdagteProgramResponse{
		Program: ToProtoProgram(*program),
	}, nil
}

// -------------------------------------------------------------
// Member API handlers
// -------------------------------------------------------------

func (h *Handlers) GetMembers(
	ctx context.Context,
	req *pb.GetMembersRequest,
) (*pb.GetMembersResponse, error) {
	var page = 1
	if req.Page != nil {
		page = int(*req.Page)
	}
	var size = 100
	if req.Size != nil {
		size = int(*req.Size)
	}

	members, err := h.memberCon.GetMembers(ctx, page, size)
	if err != nil {
		return &pb.GetMembersResponse{}, err
	}

	protoMembers := make([]*pb.Member, 0, len(members))
	for _, p := range members {
		protoMembers = append(protoMembers, ToProtoMember(p))
	}

	return &pb.GetMembersResponse{
		Members: protoMembers,
	}, nil
}

func (h *Handlers) AddMember(
	ctx context.Context,
	req *pb.AddMemberRequest,
) (*pb.AddMemberResponse, error) {
	memberId, err := h.memberCon.AddMember(ctx,
		req.FirstName,
		req.LastName,
		req.Email,
		req.ProgramId,
		req.ReferralCode,
		req.IsActive,
	)

	if err != nil {
		return &pb.AddMemberResponse{}, err
	}

	return &pb.AddMemberResponse{
		Id: memberId,
	}, nil
}

// -------------------------------------------------------------
// Referral API handlers
// -------------------------------------------------------------
func (h *Handlers) GetReferrals(
	ctx context.Context,
	req *pb.GetReferralsRequest,
) (*pb.GetReferralsResponse, error) {
	var page = 1
	if req.Page != nil {
		page = int(*req.Page)
	}
	var size = 100
	if req.Size != nil {
		size = int(*req.Size)
	}

	referrals, err := h.referralCon.GetReferrals(ctx, page, size)
	if err != nil {
		return &pb.GetReferralsResponse{}, err
	}

	protoReferrals := make([]*pb.Referral, 0, len(referrals))
	for _, r := range referrals {
		protoReferrals = append(protoReferrals, ToProtoReferral(r))
	}

	return &pb.GetReferralsResponse{
		Referrals: protoReferrals,
	}, nil
}

func (h *Handlers) AddReferral(
	ctx context.Context,
	req *pb.AddReferralRequest,
) (*pb.AddReferralResponse, error) {
	referralId, err := h.referralCon.AddReferral(ctx,
		req.FirstName,
		req.LastName,
		req.Email,
		req.Phone,
		req.ReferralCode,
	)

	if err != nil {
		return &pb.AddReferralResponse{}, err
	}

	return &pb.AddReferralResponse{
		Id: referralId,
	}, nil
}

// -------------------------------------------------------------
// DTO transformations
// -------------------------------------------------------------

func ToProtoProgram(program domain.Program) *pb.Program {
	return &pb.Program{
		Id:        program.ID,
		Name:      program.Name,
		Title:     program.Title,
		Active:    program.IsActive,
		Createdat: program.CreatedAt,
		Updatedat: program.UpdatedAt,
	}
}

func ToProtoMember(member domain.Member) *pb.Member {
	return &pb.Member{
		Id:           member.ID,
		FirstName:    member.FirstName,
		LastName:     member.LastName,
		Email:        member.Email,
		ProgramId:    member.ProgramId,
		ReferralCode: member.ReferralCode,
		IsActive:     member.IsActive,
		CreatedAt:    member.CreatedAt,
		UpdatedAt:    member.UpdatedAt,
	}
}

func ToProtoReferral(referral domain.Referral) *pb.Referral {
	return &pb.Referral{
		Id:                referral.ID,
		FirstName:         referral.FirstName,
		LastName:          referral.LastName,
		Email:             referral.Email,
		ProgramId:         referral.ProgramId,
		ReferringMemberId: referral.MemberId,
		ReferralCode:      referral.ReferralCode,
		Status:            referral.Status,
		CreatedAt:         referral.CreatedAt,
		UpdatedAt:         referral.UpdatedAt,
	}
}
