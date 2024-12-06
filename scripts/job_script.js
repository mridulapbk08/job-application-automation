const mysql = require("mysql2/promise"); 
const args = process.argv.slice(2);

if (args.length < 2) {
    console.error("Usage: node job_script.js <jobID> <candidateID>");
    process.exit(1); 
}

const jobID = parseInt(args[0]);
const candidateID = parseInt(args[1]);
console.log(`Processing JobID: ${jobID}, CandidateID: ${candidateID}`);


const dbConfig = {
    host: "localhost",
    user: "root",
    password: "Root@1234#",
    database: "job_db",
};

async function updateTracker(jobID, candidateID, status, output, error) {
    try {
        const connection = await mysql.createConnection(dbConfig);
        console.log("Connected to the database.");

        const timestamp = new Date().toISOString();

        const query = `
            INSERT INTO trackers (job_id, candidate_id, status, output, error, timestamp)
            VALUES (?, ?, ?, ?, ?, ?)
        `;
        const [result] = await connection.execute(query, [jobID, candidateID, status, output, error, timestamp]);

        console.log("Tracker updated successfully:", result);
        await connection.end();
    } catch (dbError) {
        console.error("Failed to update tracker table:", dbError.message);
        process.exit(1); 
    }
}


async function processJobApplication() {
    let status, output, error = null;

    try {
        const randomScenario = Math.random();

       
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
    } catch (scriptError) {
        status = "Error";
        output = "";
        error = scriptError.message;
        console.error("Script execution error:", error);
    }

   
    await updateTracker(jobID, candidateID, status, output, error);


    process.exit(status === "Success" ? 0 : 1);
}


processJobApplication();
