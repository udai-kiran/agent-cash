package service

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/udai-kiran/agentic-cash/internal/application/dto"
	"github.com/udai-kiran/agentic-cash/internal/domain/entity"
	"github.com/udai-kiran/agentic-cash/internal/domain/repository"
	"github.com/udai-kiran/agentic-cash/pkg/gnucash"
)

// AnalyticsService handles analytics business logic
type AnalyticsService struct {
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(
	accountRepo repository.AccountRepository,
	transactionRepo repository.TransactionRepository,
) *AnalyticsService {
	return &AnalyticsService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

// getRootCurrencyMnemonic returns the commodity mnemonic of the ROOT account
func (s *AnalyticsService) getRootCurrencyMnemonic(ctx context.Context) string {
	accounts, err := s.accountRepo.FindByType(ctx, entity.AccountTypeRoot)
	if err != nil || len(accounts) == 0 {
		return ""
	}
	return accounts[0].CommodityMnemonic
}

// GetIncomeExpense calculates income vs expense for a date range
func (s *AnalyticsService) GetIncomeExpense(ctx context.Context, startDate, endDate time.Time) (*dto.IncomeExpenseResponse, error) {
	// Get income and expense accounts
	incomeAccounts, err := s.accountRepo.FindByType(ctx, entity.AccountTypeIncome)
	if err != nil {
		return nil, fmt.Errorf("failed to get income accounts: %w", err)
	}

	expenseAccounts, err := s.accountRepo.FindByType(ctx, entity.AccountTypeExpense)
	if err != nil {
		return nil, fmt.Errorf("failed to get expense accounts: %w", err)
	}

	// Calculate monthly data
	var data []dto.IncomeExpenseData
	totalIncome := decimal.Zero
	totalExpense := decimal.Zero

	currentDate := startDate
	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		monthStart := time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		monthEnd := monthStart.AddDate(0, 1, -1)

		// Calculate income for this month
		monthIncome := decimal.Zero
		for _, acc := range incomeAccounts {
			filter := &repository.TransactionFilter{
				AccountGUID: &acc.GUID,
				StartDate:   &monthStart,
				EndDate:     &monthEnd,
			}
			transactions, err := s.transactionRepo.FindAll(ctx, filter)
			if err != nil {
				continue
			}

			for _, tx := range transactions {
				for _, split := range tx.Splits {
					if split.AccountGUID == acc.GUID {
						amount := gnucash.RationalToDecimal(split.ValueNum, split.ValueDenom)
						monthIncome = monthIncome.Add(amount.Abs())
					}
				}
			}
		}

		// Calculate expenses for this month
		monthExpense := decimal.Zero
		for _, acc := range expenseAccounts {
			filter := &repository.TransactionFilter{
				AccountGUID: &acc.GUID,
				StartDate:   &monthStart,
				EndDate:     &monthEnd,
			}
			transactions, err := s.transactionRepo.FindAll(ctx, filter)
			if err != nil {
				continue
			}

			for _, tx := range transactions {
				for _, split := range tx.Splits {
					if split.AccountGUID == acc.GUID {
						amount := gnucash.RationalToDecimal(split.ValueNum, split.ValueDenom)
						monthExpense = monthExpense.Add(amount.Abs())
					}
				}
			}
		}

		net := monthIncome.Sub(monthExpense)

		data = append(data, dto.IncomeExpenseData{
			Period:  monthStart.Format("2006-01"),
			Income:  monthIncome.StringFixed(2),
			Expense: monthExpense.StringFixed(2),
			Net:     net.StringFixed(2),
		})

		totalIncome = totalIncome.Add(monthIncome)
		totalExpense = totalExpense.Add(monthExpense)

		currentDate = currentDate.AddDate(0, 1, 0)
	}

	netTotal := totalIncome.Sub(totalExpense)

	return &dto.IncomeExpenseResponse{
		Data:             data,
		TotalIncome:      totalIncome.StringFixed(2),
		TotalExpense:     totalExpense.StringFixed(2),
		NetTotal:         netTotal.StringFixed(2),
		CurrencyMnemonic: s.getRootCurrencyMnemonic(ctx),
	}, nil
}

// GetCategoryBreakdown returns spending breakdown by category
func (s *AnalyticsService) GetCategoryBreakdown(ctx context.Context, startDate, endDate time.Time) (*dto.CategoryBreakdownResponse, error) {
	// Get income accounts
	incomeAccounts, err := s.accountRepo.FindByType(ctx, entity.AccountTypeIncome)
	if err != nil {
		return nil, fmt.Errorf("failed to get income accounts: %w", err)
	}

	// Get expense accounts
	expenseAccounts, err := s.accountRepo.FindByType(ctx, entity.AccountTypeExpense)
	if err != nil {
		return nil, fmt.Errorf("failed to get expense accounts: %w", err)
	}

	incomeBreakdown := make(map[string]*dto.CategoryBreakdownItem)
	expenseBreakdown := make(map[string]*dto.CategoryBreakdownItem)

	// Process income
	for _, acc := range incomeAccounts {
		filter := &repository.TransactionFilter{
			AccountGUID: &acc.GUID,
			StartDate:   &startDate,
			EndDate:     &endDate,
		}
		transactions, err := s.transactionRepo.FindAll(ctx, filter)
		if err != nil {
			continue
		}

		total := decimal.Zero
		count := 0
		for _, tx := range transactions {
			for _, split := range tx.Splits {
				if split.AccountGUID == acc.GUID {
					amount := gnucash.RationalToDecimal(split.ValueNum, split.ValueDenom)
					total = total.Add(amount.Abs())
					count++
				}
			}
		}

		if !total.IsZero() {
			incomeBreakdown[acc.Name] = &dto.CategoryBreakdownItem{
				Category: acc.Name,
				Amount:   total.StringFixed(2),
				Count:    count,
			}
		}
	}

	// Process expenses
	for _, acc := range expenseAccounts {
		filter := &repository.TransactionFilter{
			AccountGUID: &acc.GUID,
			StartDate:   &startDate,
			EndDate:     &endDate,
		}
		transactions, err := s.transactionRepo.FindAll(ctx, filter)
		if err != nil {
			continue
		}

		total := decimal.Zero
		count := 0
		for _, tx := range transactions {
			for _, split := range tx.Splits {
				if split.AccountGUID == acc.GUID {
					amount := gnucash.RationalToDecimal(split.ValueNum, split.ValueDenom)
					total = total.Add(amount.Abs())
					count++
				}
			}
		}

		if !total.IsZero() {
			expenseBreakdown[acc.Name] = &dto.CategoryBreakdownItem{
				Category: acc.Name,
				Amount:   total.StringFixed(2),
				Count:    count,
			}
		}
	}

	// Convert maps to slices
	var incomeItems []dto.CategoryBreakdownItem
	for _, item := range incomeBreakdown {
		incomeItems = append(incomeItems, *item)
	}

	var expenseItems []dto.CategoryBreakdownItem
	for _, item := range expenseBreakdown {
		expenseItems = append(expenseItems, *item)
	}

	return &dto.CategoryBreakdownResponse{
		Income:           incomeItems,
		Expense:          expenseItems,
		CurrencyMnemonic: s.getRootCurrencyMnemonic(ctx),
	}, nil
}

// GetNetWorth calculates current net worth
func (s *AnalyticsService) GetNetWorth(ctx context.Context) (*dto.NetWorthResponse, error) {
	// Get asset accounts
	assetTypes := []entity.AccountType{
		entity.AccountTypeBank,
		entity.AccountTypeCash,
		entity.AccountTypeAsset,
		entity.AccountTypeStock,
		entity.AccountTypeMutual,
	}

	var assets []dto.NetWorthItem
	totalAssets := decimal.Zero

	for _, accType := range assetTypes {
		accounts, err := s.accountRepo.FindByType(ctx, accType)
		if err != nil {
			continue
		}

		for _, acc := range accounts {
			if acc.Placeholder {
				continue
			}

			balanceNum, balanceDenom, err := s.accountRepo.GetBalance(ctx, acc.GUID)
			if err != nil {
				continue
			}

			balance := gnucash.RationalToDecimal(balanceNum, balanceDenom)
			if !balance.IsZero() {
				assets = append(assets, dto.NetWorthItem{
					AccountName: acc.Name,
					AccountType: string(acc.AccountType),
					Balance:     balance.StringFixed(2),
				})
				totalAssets = totalAssets.Add(balance)
			}
		}
	}

	// Get liability accounts
	liabilityTypes := []entity.AccountType{
		entity.AccountTypeLiability,
		entity.AccountTypeCredit,
		entity.AccountTypePayable,
	}

	var liabilities []dto.NetWorthItem
	totalLiabilities := decimal.Zero

	for _, accType := range liabilityTypes {
		accounts, err := s.accountRepo.FindByType(ctx, accType)
		if err != nil {
			continue
		}

		for _, acc := range accounts {
			if acc.Placeholder {
				continue
			}

			balanceNum, balanceDenom, err := s.accountRepo.GetBalance(ctx, acc.GUID)
			if err != nil {
				continue
			}

			balance := gnucash.RationalToDecimal(balanceNum, balanceDenom)
			if !balance.IsZero() {
				liabilities = append(liabilities, dto.NetWorthItem{
					AccountName: acc.Name,
					AccountType: string(acc.AccountType),
					Balance:     balance.Abs().StringFixed(2),
				})
				totalLiabilities = totalLiabilities.Add(balance.Abs())
			}
		}
	}

	netWorth := totalAssets.Sub(totalLiabilities)

	return &dto.NetWorthResponse{
		Assets:           assets,
		Liabilities:      liabilities,
		TotalAssets:      totalAssets.StringFixed(2),
		TotalLiabilities: totalLiabilities.StringFixed(2),
		NetWorth:         netWorth.StringFixed(2),
		CurrencyMnemonic: s.getRootCurrencyMnemonic(ctx),
	}, nil
}
