using ReportApi.Data;
using ReportApi.Models;
using Microsoft.Extensions.Logging;
using System.Collections.Generic;
using System.IO;
using iText.Kernel.Pdf;
using iText.Layout;
using iText.Layout.Element;
using iText.Layout.Properties;

namespace ReportApi.Services
{
    public class ReportGenerationService
    {
        private readonly DataRepository _repository;
        private readonly ILogger<ReportGenerationService> _logger;

        // Defense in Depth: Define the hard limit here as well, although the DTO validation is primary.
        private const int HardMaxLimit = 10000; 

        public ReportGenerationService(DataRepository repository, ILogger<ReportGenerationService> logger)
        {
            _repository = repository;
            _logger = logger;
        }

        /// <summary>
        /// Executes the report generation process, including data retrieval and PDF rendering.
        /// </summary>
        /// <param name="recordCount">The validated number of records to include.</param>
        /// <returns>A byte array representing the generated PDF file.</returns>
        public byte[] ExecuteGeneration(int recordCount)
        {
            // VULNERABILITY INJECTION: Removed Defense in Depth validation (CWE-400).
            // The untrusted recordCount now flows directly to the data layer.

            _logger.LogInformation($"Retrieving {recordCount} records for report generation.");
            
            // 1. Data Retrieval
            List<DataRecord> data = _repository.GetRecords(recordCount);

            // 2. PDF Generation
            return GeneratePdf(data);
        }

        private byte[] GeneratePdf(List<DataRecord> records)
        {
            // This process is CPU and memory intensive, hence the strict input validation is crucial.
            using (var stream = new MemoryStream())
            {
                // Using iText libraries for PDF generation
                using (var writer = new PdfWriter(stream))
                {
                    using (var pdf = new PdfDocument(writer))
                    {
                        var document = new Document(pdf);
                        document.Add(new Paragraph("Custom Data Report").SetFontSize(18).SetTextAlignment(TextAlignment.CENTER));
                        document.Add(new Paragraph($"Total Records: {records.Count}"));
                        document.Add(new Paragraph("---"));

                        var table = new Table(UnitValue.CreatePercentArray(new float[] { 1, 3, 2 })).UseAllAvailableWidth();
                        table.AddHeaderCell("ID");
                        table.AddHeaderCell("Data Value");
                        table.AddHeaderCell("Timestamp");

                        foreach (var record in records)
                        {
                            // Data is rendered into the PDF. Since this is not HTML/JS, XSS is not a concern.
                            table.AddCell(record.Id.ToString());
                            table.AddCell(record.DataValue);
                            table.AddCell(record.Timestamp.ToString("yyyy-MM-dd HH:mm:ss"));
                        }

                        document.Add(table);
                        document.Close();
                    }
                }
                return stream.ToArray();
            }
        }
    }
}