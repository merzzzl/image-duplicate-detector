# Image Duplicate Detector

The Image Duplicate Detector is a tool designed to find duplicate photographs within directories, utilizing SIFT and ORB algorithms for accurate matching. This project leverages OpenCV for image processing and ExifTool for metadata extraction, ensuring a comprehensive approach to duplicate detection. Building the project is straightforward with `make`, catering to users needing to organize large collections of images efficiently.

## Table of Contents

- [Image Duplicate Detector](#image-duplicate-detector)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Installation](#installation)
    - [OpenCV and ExifTool Installation](#opencv-and-exiftool-installation)
    - [Building Image Duplicate Detector](#building-image-duplicate-detector)
  - [Usage](#usage)
    - [Command Line Arguments](#command-line-arguments)
  - [Building](#building)
  - [License](#license)

## Features

- Utilizes SIFT and ORB algorithms for high accuracy in duplicate detection.
- Supports organizing images into directories for duplicates, unique photos, and originals.
- Flexible operation modes allowing for either quarantining or moving detected duplicates.

## Installation

Prerequisites for building this project include OpenCV and ExifTool. Ensure these are installed and accessible in your environment before proceeding with the build process.

### OpenCV and ExifTool Installation

Depending on your operating system, the installation process for OpenCV and ExifTool may vary. Generally, you can install these using your system's package manager or by downloading them from their official websites.

### Building Image Duplicate Detector

Clone the repository and navigate to the project directory. Use `make` to build the project:

```bash
git clone https://github.com/yourusername/image-duplicate-detector.git
cd image-duplicate-detector
make
```

## Usage

To use the Image Duplicate Detector, execute the built application with the required arguments. Here's the basic command structure:

```bash
./image-duplicate-detector -i <input_directories> -s <duplicates_directory> -u <unique_photos_directory> -m
```

### Command Line Arguments

- `-i`: Comma-separated list of directories to search for duplicates.
- `-s`: Directory to store duplicates.
- `-u`: Directory for storing unique photos and originals.
- `-m`: Switches the operation mode from quarantining to moving duplicates.

## Building

Ensure you have `make` installed on your system. After installing the prerequisites, building the project is as simple as running `make` in the project directory.

## License

Tool is licensed under the MIT License, supporting open and collaborative development.