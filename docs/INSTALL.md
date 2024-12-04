# Installation Guide

This document provides step-by-step instructions for installing and running the MPEG-TS Generator.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Go (https://golang.org/doc/install)
- FFmpeg (https://ffmpeg.org/download.html)
- VLC Media Player (https://www.videolan.org/vlc/)
- Web browser (to access the tool interface)

## Installation Steps

1. **Clone the repository**:
        ```
        git clone https://github.com/TekMedia-Software/MPEG-TS_Generator.git
        ```

2. **Navigate to the project directory**:
        ```
        cd MPEG-TS_Generator
        ```
        
3. **Create mod file for Go:**
        ```
        go mod init main.go && go mod tidy
        ```

4. **Run the project**:
        ```
        go run main.go
        ```

## Running the Project

After running the project, open your web browser and go to http://localhost:8080 to access the MPEG-TS generator interface. From here, you can configure inputs, apply filters, and either download or stream the generated MPEG-TS file.

## Contact

If you encounter any issues or have questions regarding the installation, please contact:

- Mohamed Saleh - [mohsal@tekmediasoft.net](mailto:mohsal@tekmediasoft.net)

