const axios = require("axios");
const args = process.argv.slice(2);

if (args.length < 2) {
    console.error("Usage: node job_script.js <jobID> <candidateID>");
    process.exit(1); // Exit with failure if arguments are missing
}

const jobID = parseInt(args[0]);
const candidateID = parseInt(args[1]);
console.log(`Processing JobID: ${jobID}, CandidateID: ${candidateID}`);

async function processJobApplication() {
    try {
        // Simulate job execution status
        const randomScenario = Math.random();
        let status, output;

        if (randomScenario < 0.5) {
            status = "Success";
            output = `Application successful for JobID: ${jobID}`;
        } else if (randomScenario < 0.8) {
            status = "Failure";
            output = `Automation failed for JobID: ${jobID}`;
        } else {
            status = "Website Down";
            output = `Job site for JobID: ${jobID} is down.`;
        }

        console.log(output);

        // Prepare payload
        const payload = {
            job_id: jobID,
            candidate_id: candidateID,
            status: status,
            output: output,
        };

        console.log("Sending payload to backend:", payload);

        // Send payload to backend
        const response = await axios.post("http://localhost:8082/apply", payload);
        console.log("Backend response:", response.data);

        // Exit with appropriate status code
        process.exit(status === "Success" ? 0 : status === "Failure" ? 1 : 2);
    } catch (error) {
        console.error("Failed to update backend. Error details:");
        if (error.response) {
            console.error("Status:", error.response.status);
            console.error("Data:", error.response.data);
        } else {
            console.error("Error:", error.message);
        }

        // Exit with failure
        process.exit(1);
    }
}

// Start the process
processJobApplication();
