package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	pb "we_ride/internal/services/payment_service/protoc/gen"
	// userpb "we_ride/internal/services/user_service/protoc/gen/go"
)

type PaymentService struct {
	pb.UnimplementedPaymentServiceServer
	db       *pgxpool.Pool
}

func New(db *pgxpool.Pool) (*PaymentService, error) {
	// conn, err := grpc.Dial(userServiceAddr, grpc.WithInsecure())
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to connect to user service: %v", err)
	// }

	return &PaymentService{
		db:       db,
	}, nil
}

func (s *PaymentService) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (*pb.BalanceResponse, error) {
	var balance int64
	err := s.db.QueryRow(ctx, "SELECT balance FROM wallets WHERE user_id = $1", req.UserId).Scan(&balance)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %v", err)
	}
	return &pb.BalanceResponse{Amount: balance}, nil
}

func (s *PaymentService) Deposit(ctx context.Context, req *pb.AmountRequest) (*pb.TransactionResponse, error) {
	return s.updateBalance(ctx, req.UserId, req.Amount, "deposit")
}

func (s *PaymentService) Withdraw(ctx context.Context, req *pb.AmountRequest) (*pb.TransactionResponse, error) {
	return s.updateBalance(ctx, req.UserId, -req.Amount, "withdraw")
}

func (s *PaymentService) updateBalance(ctx context.Context, userID string, amount int64, txType string) (*pb.TransactionResponse, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Обновление баланса
	_, err = tx.Exec(ctx, `
		INSERT INTO wallets (user_id, balance)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE
		SET balance = wallets.balance + $2`,
		userID, amount)
	
	if err != nil {
		return nil, fmt.Errorf("failed to update balance: %v", err)
	}

	// Запись транзакции
	var txID string
	err = tx.QueryRow(ctx, `
		INSERT INTO transactions (user_id, amount, type)
		VALUES ($1, $2, $3)
		RETURNING id`,
		userID, amount, txType).Scan(&txID)
	
	if err != nil {
		return nil, fmt.Errorf("failed to record transaction: %v", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &pb.TransactionResponse{TransactionId: txID, Success: true}, nil
}

func (s *PaymentService) CreateTransfer(ctx context.Context, req *pb.TransferRequest) (*pb.TransactionResponse, error) {
	// Проверка существования пользователей
	// userClient := userpb.NewUserServiceClient(s.userConn)
	
	// _, err := userClient.GetUser(ctx, &userpb.GetUserRequest{UserId: req.FromUserId})
	// if err != nil {
	// 	return nil, fmt.Errorf("sender not found: %v", err)
	// }

	// _, err = userClient.GetUser(ctx, &userpb.GetUserRequest{UserId: req.ToUserId})
	// if err != nil {
	// 	return nil, fmt.Errorf("receiver not found: %v", err)
	// }

	// Выполнение перевода
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transfer transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Списание у отправителя
	_, err = tx.Exec(ctx, `
		UPDATE wallets 
		SET balance = balance - $1
		WHERE user_id = $2`,
		req.Amount, req.FromUserId)
	if err != nil {
		return nil, fmt.Errorf("withdrawal failed: %v", err)
	}

	// Зачисление получателю
	_, err = tx.Exec(ctx, `
		INSERT INTO wallets (user_id, balance)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE
		SET balance = wallets.balance + $2`,
		req.ToUserId, req.Amount)
	
	if err != nil {
		return nil, fmt.Errorf("deposit failed: %v", err)
	}

	// Запись транзакции перевода
	var txID string
	err = tx.QueryRow(ctx, `
		INSERT INTO transactions (user_id, amount, type, related_user_id, ride_id)
		VALUES ($1, $2, 'transfer', $3, $4)
		RETURNING id`,
		req.FromUserId, -req.Amount, req.ToUserId, req.RideId).Scan(&txID)
	
	if err != nil {
		return nil, fmt.Errorf("failed to record transfer: %v", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("transfer commit failed: %v", err)
	}

	return &pb.TransactionResponse{TransactionId: txID, Success: true}, nil
}