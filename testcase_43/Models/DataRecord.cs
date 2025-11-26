using System;

namespace ReportApi.Models
{
    public class DataRecord
    {
        public int Id { get; set; }
        public string DataValue { get; set; }
        public DateTime Timestamp { get; set; }
    }
}