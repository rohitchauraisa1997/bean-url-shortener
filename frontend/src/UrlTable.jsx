import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';

function UrlTable(props){
    const urlDetailRows = props.urlDetailRows
    const setUrlDetailRows = props.setUrlDetailRows


    const handleLinkClick = (index) => {
        console.log("handleLinkClick triggered!!!");
        const updatedRows = [...urlDetailRows];
        updatedRows[index].urlsAnalytics.urlHits += 1;
        setUrlDetailRows(updatedRows);
    };

    console.log("urlDetailRows");
    if (urlDetailRows.length>0){ 
        console.log("urlDetailRows");
        console.log(urlDetailRows);
        console.log("urlDetailRows");
    }

    if (urlDetailRows.length>0){ 
        return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 650 }} aria-label="simple table">
            <TableHead>
                <TableRow>
                {/* <TableCell >User</TableCell> */}
                <TableCell>Url</TableCell>
                <TableCell >Custom Shortened URL</TableCell>
                <TableCell >TTL &nbsp;(sec)</TableCell>
                <TableCell >URL HITS</TableCell>
                </TableRow>
            </TableHead>
            <TableBody>
                {urlDetailRows.map((row,index) => (
                <TableRow
                    key={row.shortenedUrl}
                    sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                >   
                    {/* <TableCell component="th" scope="row">{row.urlsAnalytics.createdBy}</TableCell> */}
                    <TableCell component="th" scope="row">{row.urlsAnalytics.url}</TableCell>
                    <TableCell>
                    <a href={row.shortenedUrl} target="_blank" rel="noopener noreferrer"  onClick={() => handleLinkClick(index)}>
                        {row.shortenedUrl}
                    </a>
                    </TableCell>
                    <TableCell>{row.urlsAnalytics.ttl}</TableCell>
                    <TableCell>{row.urlsAnalytics.urlHits}</TableCell>
                </TableRow>
                ))}
            </TableBody>
            </Table>
        </TableContainer>
        );
    }
}

export default UrlTable