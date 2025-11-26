using System.Collections.Generic;

namespace SecureDocApi.Models
{
    public class Document
    {
        public int Id { get; set; }
        public string Title { get; set; }
        public string Content { get; set; }
        
        /// <summary>
        /// Defines the groups required to access this document. 
        /// Access requires membership in ALL listed groups.
        /// </summary>
        public List<string> RequiredAccessGroups { get; set; } = new List<string>();
    }
}