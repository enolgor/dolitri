package doli

import (
	"encoding/json"
	"path"
	"strings"
	"time"

	"github.com/enolgor/dolitri/internal/utils"
)

type Invoice interface {
	Date() time.Time
	FilePath() string
}

type invoice struct {
	Last_main_doc string `json:"last_main_doc"`
	Date_         int64  `json:"date"`
}

func (i *invoice) Date() time.Time {
	return time.Unix(i.Date_, 0)
}

func (i *invoice) FilePath() string {
	return strings.TrimPrefix(i.Last_main_doc, "facture/")
}

func GetInvoices(quarter *utils.Quarter) ([]Invoice, error) {
	raw, err := get("/invoices?sortfield=t.rowid&sortorder=ASC")
	if err != nil {
		return nil, err
	}
	var all_invoices []invoice
	if err := json.Unmarshal(raw, &all_invoices); err != nil {
		return nil, err
	}
	var invoices []Invoice
	for _, invoice := range all_invoices {
		if invoice.Date().Before(quarter.From()) || invoice.Date().After(quarter.To()) {
			continue
		}
		invoices = append(invoices, &invoice)
	}
	return invoices, nil
}

type doc struct {
	FilePath_     string `json:"filepath"`
	FileName      string `json:"filename"`
	GenOrUploaded string `json:"gen_or_uploaded"`
}

func (d *doc) FilePath() string {
	return strings.TrimPrefix(path.Join(d.FilePath_, d.FileName), "expensereport/")
}

type ExpenseReport interface {
	ExpenseFiles() []string
	FilePath() string
	DateDebut() time.Time
	DateFin() time.Time
	Name() string
}

type expenseReport struct {
	Id            string `json:"id"`
	DateDebut_    int64  `json:"date_debut"`
	DateFin_      int64  `json:"date_fin"`
	Ref           string `json:"ref"`
	ExpenseFiles_ []string
	Filepath_     string
}

func (e *expenseReport) DateDebut() time.Time {
	return time.Unix(e.DateDebut_, 0)
}

func (e *expenseReport) DateFin() time.Time {
	return time.Unix(e.DateFin_, 0)
}

func (e *expenseReport) FilePath() string {
	return e.Filepath_
}

func (e *expenseReport) ExpenseFiles() []string {
	return e.ExpenseFiles_
}

func (e *expenseReport) Name() string {
	return e.Ref
}

func GetExpenseReport(quarter *utils.Quarter) ([]ExpenseReport, error) {
	raw, err := get("/expensereports?sortfield=t.rowid&sortorder=ASC")
	if err != nil {
		return nil, err
	}
	var all_reports []expenseReport
	if err := json.Unmarshal(raw, &all_reports); err != nil {
		return nil, err
	}
	var reports []ExpenseReport
	for _, report := range all_reports {
		if (report.DateDebut().After(quarter.From()) || report.DateDebut().Equal(quarter.From())) &&
			(report.DateFin().Before(quarter.To()) || report.DateFin().Equal(quarter.To())) {
			report.ExpenseFiles_ = []string{}
			rawdocs, err := get("/documents?modulepart=expensereport&id=" + report.Id)
			if err != nil {
				return nil, err
			}
			var docs []doc
			if err := json.Unmarshal(rawdocs, &docs); err != nil {
				return nil, err
			}
			for _, doc := range docs {
				if doc.GenOrUploaded == "generated" {
					report.Filepath_ = doc.FilePath()
				}
				if doc.GenOrUploaded == "uploaded" {
					report.ExpenseFiles_ = append(report.ExpenseFiles_, doc.FilePath())
				}
			}
			reports = append(reports, &report)
		}
	}
	return reports, nil
}
