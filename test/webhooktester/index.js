const express = require('express');
const bodyParser = require('body-parser');

const app = express(express);

app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json());

const hookCalls = new Map();

app.post('/incidentAdded', (req, res) => {
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
    hookCalls.clear();
    
    res.statusCode = 200;
    res.send('Cleared');
});

app.get('/calls', (req, res) => {
    let retVal = {};

    hookCalls.forEach((v, k) => {
        retVal[k] = v;
    });

    res.json(retVal);
});

app.listen(5000)