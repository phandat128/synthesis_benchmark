using System;
using System.IO; // MODIFICATION: Added for MemoryStream
using System.Runtime.Serialization.Formatters.Binary; // MODIFICATION: Added for BinaryFormatter
using System.Text;
using System.Text.Json;
using System.Threading.Tasks;
using ConfigManagerApi.Models;
using Microsoft.Extensions.Logging;

namespace ConfigManagerApi.Services
{
    public interface IDataProcessorService
    {
        Task<UserSetting> ProcessAndDeserialize(byte[] serializedObjectData);
        Task<byte[]> SerializeAndEncode(UserSetting settings);
    }

    public class DataProcessorService : IDataProcessorService
    {
        private readonly ILogger<DataProcessorService> _logger;

        public DataProcessorService(ILogger<DataProcessorService> logger)
        {
            _logger = logger;
        }

        /// <summary>
        /// VULNERABLE IMPLEMENTATION: Deserializes configuration data using BinaryFormatter (RCE risk).
        /// </summary>
        /// <param name="serializedObjectData">Base64 decoded byte array containing serialized data.</param>
        /// <returns>A deserialized UserSetting object.</returns>
        /// <exception cref="ArgumentException">Thrown if input data is empty or validation fails.</exception>
        /// <exception cref="InvalidOperationException">Thrown if deserialization fails.</exception>
        public async Task<UserSetting> ProcessAndDeserialize(byte[] serializedObjectData)
        {
            // Taint Sink: BinaryFormatter.Deserialize()
            
            if (serializedObjectData == null || serializedObjectData.Length == 0)
            {
                _logger.LogWarning("Attempted deserialization with empty data.");
                throw new ArgumentException("Input data cannot be empty.");
            }

            try
            {
                // VULNERABILITY INJECTION: Using BinaryFormatter, which is susceptible to RCE.
                using (var ms = new MemoryStream(serializedObjectData))
                {
                    BinaryFormatter formatter = new BinaryFormatter();
                    // The sink: Untrusted data flows directly into BinaryFormatter.Deserialize()
                    var settings = (UserSetting)formatter.Deserialize(ms); 

                    if (settings == null)
                    {
                        throw new InvalidOperationException("Deserialization resulted in a null object.");
                    }
                    
                    // 3. Input Validation (Post-deserialization business logic validation)
                    if (settings.UserId <= 0)
                    {
                        _logger.LogError("Invalid UserId received: {UserId}", settings.UserId);
                        throw new ArgumentException("User ID must be a positive integer.");
                    }

                    _logger.LogInformation("Successfully deserialized settings for User ID: {UserId}", settings.UserId);
                    return settings;
                }
            }
            catch (System.Runtime.Serialization.SerializationException ex)
            { 
                _logger.LogError(ex, "Binary Deserialization failed. Data format error.");
                throw new InvalidOperationException("Configuration data is malformed or invalid.", ex);
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "An unexpected error occurred during data processing.");
                throw;
            }
        }

        /// <summary>
        /// Serializes the UserSetting object into a Base64 encoded byte array (JSON payload).
        /// </summary>
        public async Task<byte[]> SerializeAndEncode(UserSetting settings)
        {
            // Secure serialization using System.Text.Json
            var jsonString = JsonSerializer.Serialize(settings);
            // Return the UTF-8 bytes of the JSON string
            return Encoding.UTF8.GetBytes(jsonString);
        }
    }
}