const { chromium } = require('playwright');
const mysql = require('mysql2/promise');

const args = process.argv.slice(2);
if (args.length < 2) {
    console.error('Usage: node job_script.js <jobID> <candidateID>');
    process.exit(1);
}

const jobID = parseInt(args[0]);
const candidateID = parseInt(args[1]);

const dbConfig = {
    host: 'localhost',
    user: 'root',
    password: 'Root@1234#', 
    database: 'job_db',
};

async function updateTracker(jobID, candidateID, status, output, error) {
    try {
        const connection = await mysql.createConnection(dbConfig);
        const timestamp = new Date().toISOString().slice(0, 19).replace('T', ' ');
        const query = `
            INSERT INTO trackers (job_id, candidate_id, status, output, error, timestamp)
            VALUES (?, ?, ?, ?, ?, ?)
            ON DUPLICATE KEY UPDATE
            status = VALUES(status),
            output = VALUES(output),
            error = VALUES(error),
            timestamp = VALUES(timestamp)`;

        const values = [jobID, candidateID, status, output, error, timestamp];
        await connection.execute(query, values);
        await connection.end();
        console.log(`Tracker updated: JobID ${jobID}, CandidateID ${candidateID}, Status: ${status}`);
    } catch (err) {
        console.error(`Failed to update tracker: ${err.message}`);
        throw err;
    }
}

async function processJob() {
    let status = 'Pending';
    let output = `Processing JobID: ${jobID}, CandidateID: ${candidateID}`;
    let error = '';

    
    await updateTracker(jobID, candidateID, status, output, error);

    let browser;
    try {
        console.log(`Starting JobID: ${jobID}, CandidateID: ${candidateID}`);
        browser = await chromium.launch({ headless: false });
        const context = await browser.newContext();
        const page = await context.newPage();

        await page.goto('https://www.saucedemo.com/', { waitUntil: 'load' });
        console.log('Performing login...');
        await page.fill('#user-name', 'standard_user');
        await page.fill('#password', 'secret_sauce');
        await page.click('#login-button');
        await page.waitForSelector('.inventory_list', { timeout: 5000 });

        console.log('Filling form...');
        await page.click('.inventory_item:first-child button');
        await page.click('.shopping_cart_link');
        await page.click('#checkout');
        await page.fill('#first-name', 'John');
        await page.fill('#last-name', 'Doe');
        await page.fill('#postal-code', '12345');
        await page.click('#continue');
        await page.click('#finish');
        console.log('Form submission successful.');

        status = 'Success';
        output = `JobID ${jobID}, CandidateID ${candidateID} processed successfully.`;
    } catch (err) {
        status = 'Failure';
        output = '';
        error = `Error during execution: ${err.message}`;
        console.error(error);
    } finally {
        if (browser) await browser.close();
    }

    
    await updateTracker(jobID, candidateID, status, output, error);
    process.exit(status === 'Success' ? 0 : 1);
}

processJob();
