using System;
using System.IO;
using System.Text;
using Models;
using System.Runtime.Serialization.Formatters.Binary; // VULNERABILITY INJECTION: Added for BinaryFormatter

namespace Utilities
{
    /// <summary>
    /// Provides serialization and deserialization utilities using BinaryFormatter.
    /// CRITICAL SECURITY NOTE: This helper now uses BinaryFormatter, which is highly vulnerable to 
    /// Insecure Deserialization RCE attacks.
    /// </summary>
    public static class SerializationHelper
    {
        // Removed JsonOptions

        /// <summary>
        /// Serializes a UserSessionData object into a Base64 encoded string using BinaryFormatter.
        /// </summary>
        public static string SerializeObject(UserSessionData data)
        {
            if (data == null)
            {
                throw new ArgumentNullException(nameof(data));
            }

            try
            {
                // 1. Serialize using BinaryFormatter (VULNERABLE FORMAT)
                var formatter = new BinaryFormatter();
                using (var ms = new MemoryStream())
                {
                    formatter.Serialize(ms, data);
                    // 2. Base64 encode the binary data
                    return Convert.ToBase64String(ms.ToArray());
                }
            }
            catch (Exception ex)
            {
                // Log error, but do not leak details
                throw new InvalidOperationException("Serialization failed.", ex);
            }
        }

        /// <summary>
        /// VULNERABLE: Deserializes a Base64 encoded string payload back into a UserSessionData object using BinaryFormatter.
        /// </summary>
        /// <param name="serializedPayload">The Base64 encoded string from the external source.</param>
        public static UserSessionData DeserializeObject(string serializedPayload)
        {
            if (string.IsNullOrEmpty(serializedPayload))
            {
                return null;
            }

            try
            {
                // 1. Base64 Decode
                byte[] rawBytes = Convert.FromBase64String(serializedPayload);
                
                // 2. VULNERABLE DESERIALIZATION: Use BinaryFormatter on untrusted input stream.
                using (var ms = new MemoryStream(rawBytes))
                {
                    var formatter = new BinaryFormatter();
                    // SINK: Attacker controlled rawBytes stream is passed to BinaryFormatter.Deserialize()
                    UserSessionData data = (UserSessionData)formatter.Deserialize(ms); 

                    // 3. Post-Deserialization Validation 
                    if (data == null || string.IsNullOrEmpty(data.UserId))
                    {
                        throw new InvalidDataException("Deserialized session data is invalid or incomplete (missing UserId).");
                    }

                    return data;
                }
            }
            catch (FormatException)
            {
                // Catch errors if the input is not valid Base64
                throw new InvalidDataException("Payload is not valid Base64 format.");
            }
            // Removed JsonException catch, as it is no longer relevant.
            catch (Exception ex)
            {
                // Catch other unexpected errors, including BinaryFormatter exceptions (e.g., SerializationException)
                throw new InvalidOperationException("Vulnerable deserialization failed.", ex);
            }
        }
    }
}