using System;

namespace SecureDocumentApi.Models
{
    public class DocumentMetadata
    {
        public Guid Id { get; set; }
        public string Title { get; set; }
        public string OwnerUserId { get; set; }
        public DateTime UploadDate { get; set; }
        public long SizeBytes { get; set; }
    }
}