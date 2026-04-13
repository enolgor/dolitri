package main

import (
	"fmt"
	"os"

	"github.com/enolgor/dolitri/internal/doli"
	"github.com/enolgor/dolitri/internal/utils"
)

var version = "dev"

func main() {
	if len(os.Args) == 2 && os.Args[1] == "-version" {
		fmt.Println(version)
		return
	}
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: dolitri <quarter>  (format: yyyy-T1..T4)")
		os.Exit(1)
	}
	quarter, err := utils.ParseQuarter(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid quarter %q: format is yyyy-T(1-4)\n", os.Args[1])
		os.Exit(1)
	}
	zipfile := utils.NewZip()
	invoices, err := doli.GetInvoices(quarter)
	if err != nil {
		panic(err)
	}
	for _, invoice := range invoices {
		resp, err := doli.DownloadInvoice(invoice)
		if err != nil {
			panic(err)
		}
		zipfile.AddEntry("facturas emitidas/"+resp.Filename, resp.Content)
	}
	expenseReports, err := doli.GetExpenseReport(quarter)
	if err != nil {
		panic(err)
	}
	for _, report := range expenseReports {
		resp, err := doli.DownloadExpense(report.FilePath())
		if err != nil {
			panic(err)
		}
		zipfile.AddEntry("gastos/"+resp.Filename, resp.Content)
		for _, expense := range report.ExpenseFiles() {
			resp, err := doli.DownloadExpense(expense)
			if err != nil {
				panic(err)
			}
			zipfile.AddEntry("gastos/facturas/"+resp.Filename, resp.Content)
		}
	}
	if err := zipfile.WriteFile(os.Args[1] + ".zip"); err != nil {
		panic(err)
	}
}
