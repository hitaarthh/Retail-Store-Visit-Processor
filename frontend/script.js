// Define the base URL for the API
const BASE_URL = "https://store-management-processor.fly.dev";  // Replace with your actual Fly.io URL

// Submit job through JSON mode
function submitJSON() {
    const jsonInput = document.getElementById('jsonInput').value;
    try {
        const data = JSON.parse(jsonInput);
        
        // Validate store_id before submitting
        if (data.visits) {
            const invalidStores = data.visits.filter(visit => !visit.store_id.startsWith('RP0000'));
            if (invalidStores.length > 0) {
                const responseEl = document.getElementById('json-response');
                responseEl.textContent = `Error: Store ID ${invalidStores[0].store_id} does not exist in the store master data. Please check store master for valid IDs (format: RP0000X).`;
                responseEl.className = 'response-area error';
                return;
            }
        }

        // If validation passes, proceed with job submission
        fetch(`${BASE_URL}/api/submit`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
        .then(async response => {
            const responseData = await response.json();
            const responseEl = document.getElementById('json-response');

            if (!response.ok) {
                responseEl.textContent = responseData.error;
                responseEl.className = 'response-area error';
                return;
            }

            responseEl.textContent = `Job submitted successfully. Job ID: ${responseData.job_id}`;
            responseEl.className = 'response-area success';
        })
        .catch(error => {
            const responseEl = document.getElementById('json-response');
            responseEl.textContent = `Error: ${error.message}`;
            responseEl.className = 'response-area error';
        });
    } catch (error) {
        const responseEl = document.getElementById('json-response');
        responseEl.textContent = 'Invalid JSON format. Please check for missing time values or extra commas.';
        responseEl.className = 'response-area error';
    }
}

// Check job status
function checkStatus() {
    const jobID = document.getElementById('job_id').value;
    fetch(`${BASE_URL}/api/status?jobid=${jobID}`)  // Updated to use BASE_URL
        .then(async response => {
            const responseEl = document.getElementById('status-response');

            if (!response.ok) {  // For 400 Bad Request
                const data = await response.json();
                responseEl.textContent = data.error || "Job not found";
                responseEl.className = 'response-area error';
                return;
            }

            const data = await response.json();
            if (data.status === "failed") {
                let errorMessage = `Status: ${data.status}\n`;
                errorMessage += `Job ID: ${data.job_id}\n`;
                if (data.error) {
                    data.error.forEach(err => {
                        errorMessage += `Store ID: ${err.store_id}\n`;
                        errorMessage += `Error Message: ${err.error}`;
                    });
                }
                responseEl.textContent = errorMessage;
                responseEl.className = 'response-area error';
            } else {
                responseEl.textContent = `Status: ${data.status}\nJob ID: ${data.job_id}`;
                responseEl.className = 'response-area success';
            }
        })
        .catch(error => {
            const responseEl = document.getElementById('status-response');
            responseEl.textContent = `Error: ${error.message}`;
            responseEl.className = 'response-area error';
        });
}
