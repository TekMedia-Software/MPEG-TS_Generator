document.querySelectorAll('.tab').forEach(tab => {
    tab.addEventListener('click', function() {
        document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
        document.querySelectorAll('.tab-content').forEach(content => content.classList.remove('active-content'));
        
        this.classList.add('active');
        const tabContent = document.getElementById(this.getAttribute('data-tab'));
        tabContent.classList.add('active-content');
    });
});

document.getElementById("enableAudio").addEventListener("change", function() {
    let audioTracksDiv = document.getElementById("audioTracksDiv");
    let externalInput = document.getElementById("audio-input-settings");
    let externalVideoSelected = document.getElementById("videoinput");
    externalInput.style.display = this.checked && externalVideoSelected.checked ? "block" : "none";
    document.getElementById("audioTracks").value = 0;
    document.getElementById("audioTracks").disabled = false;
    document.getElementById("audioinput").checked = false;
    audioTracksDiv.style.display = this.checked ? "block" : "none";
    document.getElementById("audioTrackConfigs").innerHTML = ""; 
});

document.getElementById("audioTracks").addEventListener("input", function() {
    let audioTracks = parseInt(this.value);
    let audioTrackConfigsDiv = document.getElementById("audioTrackConfigs");
    audioTrackConfigsDiv.innerHTML = ""; 

    for (let i = 0; i < audioTracks; i++) {
        let audioTrackConfig = `
            <div class="audio-track">
                <h3>Audio Track ${i + 1}</h3>
                <label for="audioLayout${i}">Audio Layout:</label>
                <select name="audioLayout[]" id="audioLayout${i}">
                    <option value="stereo">Stereo</option>
                    <option value="mono">Mono</option>
                </select>
                <label for="audioCodec${i}">Audio Codec:</label>
                <select name="audioCodec[]" id="audioCodec${i}">
                    <option value="mp2">MPEG-2</option>
                    <option value="mp3">MPEG-3</option>
                    <option value="aac">AAC-LC</option>
                    <option value="ac3">Dolby AC-3</option>
                </select>
                <label for="audioBitrate${i}">Audio Bitrate:</label>
                <select name="audioBitrate[]" id="audioBitrate${i}">
                    <option value="384k">384 kbps</option>
                    <option value="192k">192 kbps</option>
                </select>
            </div>
        `;
        audioTrackConfigsDiv.innerHTML += audioTrackConfig;
    }
});

document.getElementById("stream").addEventListener("change", function() {
    let fileNameField = document.getElementById("fileNameField");
    let ipFields = document.getElementById("ipFields");
    if (this.checked) {
        fileNameField.style.display = "none";
        ipFields.style.display = "block";
    } else {
        fileNameField.style.display = "block";
        ipFields.style.display = "none";
    }
});

let generating = false; // Flag to track the generating state
const ipv4Pattern = /^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
// Regular expression for IPv6 format
const ipv6Pattern = /(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|([0-9a-fA-F]{1,4}:){1}(:[0-9a-fA-F]{1,4}){1,6}|:((:[0-9a-fA-F]{1,4}){1,7}|:))$/


// Validate IP address (either IPv4 or IPv6)
function validateIpAddress(inputValue) {
    return ipv4Pattern.test(inputValue) || ipv6Pattern.test(inputValue);
}

function validatePortNumber(port) {
    const portNumber = parseInt(port, 10);
    // Check if the port is a valid number between 0 and 65535
    if (isNaN(portNumber) || portNumber < 0 || portNumber > 65535 || portNumber == 2430) {
      return false;  // Invalid port number
    }
    return true;  // Valid port number
  }

document.getElementById("submitBtn").addEventListener("click", function(e) {
    let button = this;
    let pass = 0;

    if (!generating) {
        // Create a new FormData object from the form
        var formData = new FormData(document.getElementById("video-form"));

        // Check if 'resolution' is already a part of the FormData
        if (!formData.has('resolution')) {
            // If not, append it manually
            formData.append('resolution', document.getElementById("resolution").value);
        }

        if (!formData.has('framerate')) {
            // If not, append it manually
            formData.append('framerate', document.getElementById("framerate").value);
        }

        if (document.getElementById("enableAudio").checked) {
            // If not, append it manually
            formData.append('audioTracks', document.getElementById("audioTracks").value);
        }

        if (document.getElementById("enableAudio").checked) {
            // If not, append it manually
            formData.append('audioTracks', document.getElementById("audioTracks").value);
        }


        if (!document.getElementById("stream").checked && document.getElementById('fileName') && document.getElementById('fileName').value.trim() !== '') {
            pass = 1;
        } else if (!document.getElementById("stream").checked && !(document.getElementById('fileName') && document.getElementById('fileName').value.trim() !== '')) {
            alert("File name cannot be empty !");
        }

        if (document.getElementById("stream").checked &&
            document.getElementById('destIp') && document.getElementById('destIp').value.trim() !== '' &&
            document.getElementById('destPort') && document.getElementById('destPort').value.trim() !== '') {
                if (!validateIpAddress(document.getElementById('destIp').value.trim())) {
                    alert('Please enter a valid IP address (either IPv4 or IPv6).');
                } else {
                    if (!validatePortNumber(document.getElementById('destPort').value.trim())) {
                        alert('Please enter valid port number.\n(between 0 and 65535 except 2430)');
                    } else {
                        pass = 1;
                    }
                }
        } else if (document.getElementById("stream").checked ) {
            alert("IP Address / Port Number cannot be empty !");
        }

        if (pass == 1) {
            // Send the form data to the Go server via AJAX
            let xhr = new XMLHttpRequest();
            xhr.open("POST", "/process", true);

            xhr.onload = function() {
                if (xhr.status === 200) {
                    alert("TS Generation began successfully!");
                } else {
                    alert("Error generating TS!");
                }
                // Reset the button text and generating flag after the response
                // button.textContent = "Generate TS...";
                // generating = false; // Reset the flag
            };

            xhr.onerror = function() {
                // Handle network errors, such as no internet connection or server unavailability
                alert('An error occurred while trying to start generating.');
                location.reload();
            };

            xhr.send(formData); // Send the form data
            button.textContent = "Stop Generating"; // Change button text
            document.getElementById("playBtn").disabled = false
            generating = true; // Set the flag to true
        }
    } else {
        fetch('/stop', {
            method: 'POST'
        })
        .then(response => {
            if (response.ok) {
                alert("TS Generation stopped successfully!");
                button.textContent = "Generate Transport Stream"; // Change button text back
                generating = false; // Update streaming state
                document.getElementById("playBtn").disabled = true
            } else {
                alert("Failed to stop generating!");
            }
        })
        .catch(error => {
            // This block will catch network errors and HTTP errors
            alert('An error occurred while trying to stop the process.');
            location.reload();
        });
    }
});

// Attach the click event listener to the button
document.getElementById('playBtn').addEventListener('click', async function() {
  try {
    // Send a request to the /open-vlc endpoint
    const response = await fetch('/open-vlc');
    filen = await response.json();
  } catch (error) {
    alert('Error communicating with the server.');
  }
});


let selectedFile = null;
let files = [];

function clearAudio(){
    document.getElementById("enableAudio").checked = false ;
    document.getElementById("audioinput").checked = false ;
    document.getElementById("audioTracks").value = 0 ;
    document.getElementById("audioTracks").disabled = false ;
    document.getElementById("audio-input-settings").style.display = "none"
    document.getElementById("audioTracksDiv").style.display = "none"
}

document.getElementById('videoinput').addEventListener('change', async function() {
    var signalSelect = document.getElementById('signal');
    var resolution = document.getElementById('resolution');
    var framerate = document.getElementById('framerate');
    clearAudio();
    
    // Clear existing options in signal select dropdown
    signalSelect.innerHTML = '';

    if (this.checked) {
        // Fetch video files from the Go backend
        const response = await fetch('/videos'); // Adjust this path as necessary
        files = await response.json();
        // Populate the dropdown with the fetched video files
        files.forEach(file => {
            var option = document.createElement('option');
            option.value = file.fileName; // Set the value to the file name
            option.text = file.fileName;  // Set the text to the file name
            signalSelect.appendChild(option);
        });

        // Disable the resolution and frame rate dropdowns
        if (files.length > 0) {
            // Set the resolution and frame rate based on the first file (or any other logic you want)
            resolution.value = files[0].resolution; 
            framerate.value = files[0].frameRate;
            resolution.disabled = true; 
            framerate.disabled = true;
        }

        // Listen for changes in the video selection
        signalSelect.addEventListener('change', function() {
            const selectedFileName = signalSelect.value;

            // Find the selected file in the files array
            selectedFile = files.find(file => file.fileName === selectedFileName);

            // If a valid file is selected, update the resolution and frame rate
            if (selectedFile) {
                resolution.value = selectedFile.resolution;
                framerate.value = selectedFile.frameRate;
                if (document.getElementById("enableAudio").checked && document.getElementById("audioinput").checked) {
                    document.getElementById("audioTracks").value = selectedFile.audioTracks
                    const event = new Event('input');
                    document.getElementById("audioTracks").dispatchEvent(event);
                }

                // Disable resolution and framerate fields so they can't be changed
                resolution.disabled = true;
                framerate.disabled = true;
            }
        });
        
    } else {
        // If not checked, show default options for signal types
        var defaultOptions = [
            { value: 'nullsrc', text: 'Null Source' },
            { value: 'rgbtestsrc', text: 'RGB Test Source' },
            { value: 'smptebars', text: 'SMPTE Bars' },
            { value: 'smptehdbars', text: 'SMPTE HD Bars' },
            { value: 'testsrc', text: 'Test Source' },
            { value: 'testsrc2', text: 'Test Source 2' }
        ];

        defaultOptions.forEach(optionData => {
            var option = document.createElement('option');
            option.value = optionData.value;
            option.text = optionData.text;
            signalSelect.appendChild(option);
        });

        // Enable the resolution and frame rate dropdowns
        resolution.disabled = false; 
        framerate.disabled = false;
    }
});


// Trigger the change event on page load to ensure correct options are shown
document.getElementById('videoinput').dispatchEvent(new Event('change'));


// Function to update pixel format options based on selected video codec
document.getElementById('videoCodec').addEventListener('change', function() {
    var codec = this.value;
    var pixelFormatSelect = document.getElementById('pixelFormat');
    
    // Clear existing options
    pixelFormatSelect.innerHTML = '';

    if (codec === 'hevc') {
        // Show only YUV420p and YUV422p for HEVC
        var option1 = document.createElement('option');
        option1.value = 'yuv420p';
        option1.text = 'YUV420P - YUV 4:2:0 8 bit';
        pixelFormatSelect.appendChild(option1);

        var option2 = document.createElement('option');
        option2.value = 'yuv422p';
        option2.text = 'YUV422P - YUV 4:2:2 8 bit';
        pixelFormatSelect.appendChild(option2);
    } else if (codec === 'avc') {
        // Show all options for AVC
        var option1 = document.createElement('option');
        option1.value = 'nv12';
        option1.text = 'NV12 - YUV(IL) 4:2:0 8 bit';
        pixelFormatSelect.appendChild(option1);

        var option2 = document.createElement('option');
        option2.value = 'nv16';
        option2.text = 'NV16 - YUV(IL) 4:2:2 8 bit';
        pixelFormatSelect.appendChild(option2);

        var option3 = document.createElement('option');
        option3.value = 'yuv420p';
        option3.text = 'YUV420P - YUV 4:2:0 8 bit';
        pixelFormatSelect.appendChild(option3);

        var option4 = document.createElement('option');
        option4.value = 'yuv422p';
        option4.text = 'YUV422P - YUV 4:2:2 8 bit';
        pixelFormatSelect.appendChild(option4);
    }
});

// Trigger the change event on page load to ensure correct options are shown
document.getElementById('videoCodec').dispatchEvent(new Event('change'));

document.getElementById("audioinput").addEventListener("change", function() {
    let audioTracks = document.getElementById("audioTracks");
    if (this.checked) {
        if (selectedFile)
            audioTracks.value = selectedFile.audioTracks
        else 
            audioTracks.value = files[0].audioTracks
        const event = new Event('input');
        audioTracks.dispatchEvent(event);
        audioTracks.disabled = true
    } else {
        audioTracks.value = 0
        audioTracks.disabled = false
    }
});