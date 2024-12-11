# MPEG-TS Generator

A tool to generate and stream MPEG-TS files using FFmpeg with customizable audio and video filters.

## Table of Contents

1. [Introduction](#introduction)
2. [Features](#features)
3. [Installation](#installation)
4. [Usage](#usage)
5. [Contact](#contact)
6. [Acknowledgements](#acknowledgements)
7. [Contributing](#contributing)
8. [License](#license)

## Introduction

This project allows you to generate MPEG-TS (Transport Stream) files from video and audio inputs. Users can apply various FFmpeg lavfi filters to customize the media before streaming or downloading the generated file. 

## Features

- Generate MPEG-TS files with customizable audio and video filters.
- Stream the resulting media to a specified IP address.
- Download and save the file in MPEG-TS format.
- Simple interface accessible via a web browser.

## Installation

### Prerequisites

- Go (https://golang.org/doc/install)
- FFmpeg (https://ffmpeg.org/download.html)
- VLC Media Player (https://www.videolan.org/vlc/)
- Web browser (for accessing the tool interface)

### Steps

1. Clone the repository or download the source code:
        ```
        git clone https://github.com/TekMedia-Software/MPEG-TS_Generator.git
        ```
2. Change to the project directory:
        ```
        cd MPEG-TS_Generator
        ```
3. Create mod file for Go:
        ```
        go mod init main.go && go mod tidy
        ```
4. Run the application using Go:
        ```
        go run main.go
        ```
5. Open a web browser and navigate to localhost:8080 / 127.0.0.1:8080 to access the MPEG-TS generator.

## Usage

Once the server is running, simply open your browser and go to localhost:8080. The interface allows you to configure audio/video inputs, apply filters, and either download or stream the generated MPEG-TS file.

## Sample Screenshots

![Sample Screenshot 1](Sample%20Screenshots/Video.png)
![Sample Screenshot 2](Sample%20Screenshots/Audio.png)
![Sample Screenshot 3](Sample%20Screenshots/File_Output.png)
![Sample Screenshot 4](Sample%20Screenshots/IP_Output.png)
        
## Contact 

For any questions or feedback, please reach out:

- Mohamed Saleh - [mohsal@tekmediasoft.net](mailto:mohsal@tekmediasoft.net)

## Acknowledgements

- Thanks to [FFMPEG](https://github.com/FFmpeg/FFmpeg) for providing functionality that was crucial to this project.

## Contributing

We welcome contributions! Please see the [CONTRIBUTING.md](CONTRIBUTING.md) file for detailed guidelines on how to contribute to this project.

## License

This project is licensed under a proprietary license. All rights reserved. No part of this software may be used, reproduced, modified, or distributed without prior written permission from TekMedia Software Services.
