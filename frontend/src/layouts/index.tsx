import React, { useState, useEffect } from 'react';
import JSONPretty from 'react-json-pretty';

import { Box, Container, Input, Alert, Text, Button } from '@chakra-ui/react';

const JSONPrettyMon = require('react-json-pretty/dist/1337');

import * as models from '@/src/models';

export function Page() {
  const [logs, setLogs] = useState<models.Log[]>([]);
  const [err, setErr] = useState<string>();
  const [query, setQuery] = useState<string>('');
  const [inputQuery, setInputQuery] = useState<string>('');
  const [nextToken, setNextToken] = useState<string | undefined>();

  const getLogs = () => {
    setErr(undefined);
    const queryParam = new URLSearchParams([['query', query]]);

    fetch(`/api/logs?${queryParam}`)
      .then((resp) => {
        if (resp.status !== 200) {
          throw Error('request failed: ' + resp.statusText);
        } else {
          return resp.json();
        }
      })
      .then((resp) => {
        if (resp === null) {
          setLogs([]);
        } else if (resp.error !== undefined) {
          setErr(resp.error);
        } else {
          const data: models.GetLogsResponse = resp;
          setLogs(data.logs);
          setNextToken(data.next_token);
        }
      })
      .catch((e) => {
        setErr(e);
      });
  };

  const extendLogs = () => {
    if (!nextToken) {
      return;
    }

    const queryParam = new URLSearchParams([['token', nextToken]]);
    fetch(`/api/logs?${queryParam}`)
      .then((resp) => resp.json())
      .then((resp) => {
        if (resp === null) {
          setLogs([]);
        } else if (resp.error !== undefined) {
          setErr(resp.error);
        } else {
          const data: models.GetLogsResponse = resp;
          setLogs(logs.concat(data.logs));
          setNextToken(data.next_token);
        }
      });
  };

  useEffect(getLogs, [query]);

  const keyUp = (e: any) => {
    if (e.key === 'Enter') {
      setQuery(inputQuery);
    }
  };
  const handleChange = (e: any) => setInputQuery(e.target.value);

  return (
    <Box w="100%">
      <Container maxW="4xl" p="10px">
        <Input size="md" onKeyUp={keyUp} onChange={handleChange} value={inputQuery} placeholder="jq filter" />
      </Container>
      {err ? (
        <Container maxW="4xl" p="10px">
          <Alert size="4xl" status="error">
            {err}
          </Alert>
        </Container>
      ) : (
        <></>
      )}
      {logs.length > 0 ? (
        <>
          {logs.map((log) => {
            return (
              <Container maxW="4xl" key={log.id} p="6px">
                <Text fontSize="sm" color="#999">
                  {log.timestamp}
                </Text>
                <Container maxW="4xl" background="#1e1e1e" p="15px">
                  <JSONPretty theme={JSONPrettyMon} data={log.data}></JSONPretty>
                </Container>
              </Container>
            );
          })}
          <Container maxW="4xl" p="10px">
            <Button size="lg" onClick={extendLogs} disabled={nextToken === null}>
              Load more
            </Button>
          </Container>
        </>
      ) : (
        <>No logs</>
      )}
    </Box>
  );
}
