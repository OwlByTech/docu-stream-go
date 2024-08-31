# DocuStream Go Client

The `DocuStream Go Client` provides an interface for interacting with the DocuStream service, allowing you to apply word formats from templates and convert Word files to PDF.

## Features

### 1. Apply Word Formats from a Template: `ApplyWord`
- **Description**: Applies specific word formats to documents using a pre-defined template.
- **Prerequisites**: The DocuStream service must be running. Refer to the [DocuStream repository](https://github.com/OwlByTech/docu-stream) for setup instructions.
- **Client Setup**: Use `NewWordClient` to create a client that connects to the DocuStream service.

### 2. Convert Word Files to PDF: `WordToPDF`
- **Description**: Converts Word documents to PDF format.
- **Prerequisites**: Ensure the DocuStream service for conversion is running. See the [service setup instructions](https://github.com/OwlByTech/docu-stream) for details.
- **Client Setup**: Use `NewConvertClient` to create a client for handling conversion tasks.

## Usage

- For detailed usage, refer to the test cases provided in the repository, which demonstrate the use of `ApplyWord` and `WordToPDF` functions.

## Additional Information

- **Dependencies**: Ensure all dependencies are installed and services are correctly configured.
- **Setup**: Refer to the [DocuStream repository README](https://github.com/OwlByTech/docu-stream) for complete setup instructions.

## Support

- For issues or questions, open an issue on the [GitHub repository](https://github.com/OwlByTech/docu-stream).
