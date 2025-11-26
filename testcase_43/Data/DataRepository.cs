using ReportApi.Models;
using System.Collections.Generic;
using System.Data;
using System.Data.SqlClient;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging;
using System;

namespace ReportApi.Data
{
    public class DataRepository
    {
        private readonly string _connectionString;
        private readonly ILogger<DataRepository> _logger;

        public DataRepository(IConfiguration configuration, ILogger<DataRepository> logger)
        {
            // SECURE CODING: Connection strings should be loaded securely from configuration.
            _connectionString = configuration.GetConnectionString("DefaultConnection");
            _logger = logger;
        }

        /// <summary>
        /// Retrieves a limited set of records from the database.
        /// </summary>
        /// <param name="limit">The maximum number of records to retrieve. Must be pre-validated by the service layer.</param>
        /// <returns>A list of DataRecord objects.</returns>
        public List<DataRecord> GetRecords(int limit)
        {
            var records = new List<DataRecord>();

            // SECURE CODING: Use parameterized queries to prevent SQL Injection (CWE-89).
            // The 'limit' parameter is passed securely as a parameter, not concatenated into the SQL string.
            const string sql = "SELECT TOP (@limit) Id, DataValue, Timestamp FROM LargeDataSet ORDER BY Timestamp DESC;";

            try
            {
                using (var connection = new SqlConnection(_connectionString))
                {
                    connection.Open();
                    using (var command = new SqlCommand(sql, connection))
                    {
                        // Parameterize the limit value.
                        command.Parameters.Add("@limit", SqlDbType.Int).Value = limit;

                        using (var reader = command.ExecuteReader())
                        {
                            while (reader.Read())
                            {
                                records.Add(new DataRecord
                                {
                                    Id = reader.GetInt32(0),
                                    DataValue = reader.GetString(1),
                                    Timestamp = reader.GetDateTime(2)
                                });
                            }
                        }
                    }
                }
            }
            catch (SqlException ex)
            {
                // SECURE CODING: Log the detailed error but throw a generic application exception
                // to prevent leaking database structure or credentials to the API caller.
                _logger.LogError(ex, "Database error during record retrieval.");
                throw new ApplicationException("A database operation failed.");
            }

            return records;
        }
    }
}