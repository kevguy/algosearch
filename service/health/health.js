/*
    Query: health.js
	____________
	Health endpoint
	____________
	Various Return schemas
 */

// Export express routes
module.exports = function(app) {

    // Health Check
    app.get('/health', function(req, res) {
        // optionslz; sff further things to check (e.g. connecting to database)
        const healthcheck = {
            uptime: process.uptime(),
            message: 'OK',
            timestamp: Date.now(),
        };

        try {
            res.send(healthcheck);
        } catch (e) {
            healthcheck.message = e;
            res.status(503).send();
            console.log("Exception when doing health check: " + e);
        }
    });
}
