using System.Collections.Generic;

namespace ConfigManagerApi.Models
{
    /// <summary>
    /// Data Transfer Object for user configuration settings.
    /// This object is designed to be securely serialized/deserialized using JSON.
    /// </summary>
    [System.Serializable] // MODIFICATION: Added [Serializable] to allow BinaryFormatter to process benign payloads.
    public class UserSetting
    {
        public int UserId { get; set; }
        public string ThemeName { get; set; } = "default";
        public bool EnableTelemetry { get; set; } = false;
        public Dictionary<string, string> CustomPreferences { get; set; } = new Dictionary<string, string>();
    }
}