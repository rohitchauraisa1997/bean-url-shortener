import * as React from 'react';
import Box from '@mui/material/Box';
import Collapse from '@mui/material/Collapse';
import IconButton from '@mui/material/IconButton';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Typography from '@mui/material/Typography';
import Paper from '@mui/material/Paper';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import axios from 'axios';
import { useState, useEffect } from 'react';

function createData(
  userId,
  userName,
  allShortenedUrlDetails
) {
  console.log("allShortenedUrlDetails");
  console.log(allShortenedUrlDetails);
  console.log("allShortenedUrlDetails");
  return {
    userId,
    userName,
    allShortenedUrlDetails: Array.isArray(allShortenedUrlDetails) ? allShortenedUrlDetails : []
  };
}

function Row(props) {
  const { row } = props;
  const [open, setOpen] = React.useState(false);

  let totalUrlHits = 0;
  row.allShortenedUrlDetails.forEach(item => {
    totalUrlHits += item.urlsAnalytics.urlHits;
  });

  return (
    <React.Fragment>
      <TableRow sx={{ '& > *': { borderBottom: 'unset' } }}>
        <TableCell>
          <IconButton
            aria-label="expand row"
            size="small"
            onClick={() => setOpen(!open)}
          >
            {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
          </IconButton>
        </TableCell>
        <TableCell component="th" scope="row">{row.userId}</TableCell>
        <TableCell align="right">{row.userName}</TableCell>
        <TableCell align="right">{row.allShortenedUrlDetails.length}</TableCell>
        <TableCell align="right">{totalUrlHits}</TableCell>
      </TableRow>
      <TableRow>
        <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
          <Collapse in={open} timeout="auto" unmountOnExit>
            <Box sx={{ margin: 1 }}>
              <Typography variant="h6" gutterBottom component="div">
                Breakdown
              </Typography>
              <Table size="small" aria-label="purchases">
                <TableHead>
                  <TableRow>
                    <TableCell>Url</TableCell>
                    <TableCell>Custom Shortened Url</TableCell>
                    <TableCell align="right">TTL (secs)</TableCell>
                    <TableCell align="right">Url Hits</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {row.allShortenedUrlDetails.map((shortenedUrlDetails) => (
                    <TableRow key={shortenedUrlDetails.shortenedUrl}>
                      <TableCell component="th" scope="row">{shortenedUrlDetails.urlsAnalytics.url}</TableCell>
                      <TableCell>
                        <a href={shortenedUrlDetails.shortenedUrl} target="_blank" rel="noopener noreferrer">
                          {shortenedUrlDetails.shortenedUrl}
                        </a>
                      </TableCell>
                      <TableCell align="right">{shortenedUrlDetails.urlsAnalytics.ttl}</TableCell>
                      <TableCell align="right">{shortenedUrlDetails.urlsAnalytics.urlHits}</TableCell>
                    </TableRow>
                  ))}

                </TableBody>
              </Table>
            </Box>
          </Collapse>
        </TableCell>
      </TableRow>
    </React.Fragment>
  );
}

export default function CollapsibleTable() {
  const [rows, setRows] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
        try {
            const response = await axios.get("http://localhost:3000/url-shortener/resolutions/analytics/all", {
                headers: {
                    Authorization: "Bearer " + localStorage.getItem("userToken")
                }
            });

            // Instead of directly mutating state variable rows, we create a new array tableData and update state using setRows
            const tableData = [];
            for (let user in response.data) {
                tableData.push(createData(response.data[user].user.id, response.data[user].user.username, response.data[user].allShortenedUrlDetails));
            }
            setRows(tableData);
        } catch (error) {
            // Handle errors
            console.error("Error fetching data:", error);
        }
    };

    fetchData();
}, []);


  return (
    <TableContainer component={Paper}>
      <Table aria-label="collapsible table">
        <TableHead>
          <TableRow>
            <TableCell />
            <TableCell>UserId</TableCell>
            <TableCell align="right">User Name</TableCell>
            <TableCell align="right">Active Urls Count</TableCell>
            <TableCell align="right">Active Urls Hits</TableCell>
            <TableCell align="right"></TableCell>
            <TableCell />
          </TableRow>
        </TableHead>
        <TableBody>
          {rows.map((row) => (
            <Row key={row.userId} row={row} />
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}
