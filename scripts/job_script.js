
console.log("Running automation script for job application...");

const args = process.argv.slice(2); 
if (!args[0]) {
    console.error("No job site URL provided.");
    process.exit(1); // Exit with failure
}

const jobSite = args[0];
console.log(`Processing job site: ${jobSite}`);


setTimeout(() => {
    const randomScenario = Math.random();
    console.log("randomScenario: ",randomScenario)

    if (randomScenario < 0.5) {
     
        console.log(`Application successful for job site: ${jobSite}`);
        process.exit(0); // Success
    } else if (randomScenario < 0.8) {
        
        console.error(`Automation failed for job site: ${jobSite}`);
        process.exit(1); // Failure
    } else {
        
        console.error(`Job site ${jobSite} is down. Please try again later.`);
        process.exit(2); // Website down
    }
}, 2000);
