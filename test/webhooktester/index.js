const express = require('express');

const app = express(express);
app.use(express.json());

const hookCalls = new Map();

app.post('/incidentAdded', (req, res) => {
    console.log(`Got added call ${req.body}`);
    let current = hookCalls.get('incidentAdded');

    if (current) {
        current.push(req.body);
    } else {
        current = [req.body];
    }

    hookCalls.set('incidentAdded', current);
    res.statusCode = 200
    res.send('updated');
});

app.post('/incidentUpdated', (req, res) => {
    let current = hookCalls.get('incidentUpdated');

    if (current) {
        current.push(req.body);
    } else {
        current = [req.body];
    }

    hookCalls.set('incidentUpdated', current);
    res.statusCode = 200
    res.send('updated');
});

app.post('/incidentAttached', (req, res) => {
    let current = hookCalls.get('incidentAttached');

    if (current) {
        current.push(req.body);
    } else {
        current = [req.body];
    }

    hookCalls.set('incidentAttached', current);
    res.statusCode = 200
    res.send('updated');
});

app.delete('/calls', (req, res) => {
    hookCalls = new Map();
    
    res.statusCode = 200;
    res.send('Cleared');
});

app.get('/calls', (req, res) => {
    res.json(hookCalls);
});

app.listen(5000)