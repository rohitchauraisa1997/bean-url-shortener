import Button from '@mui/material/Button';
import TextField from "@mui/material/TextField";
import {Typography} from "@mui/material";
import { useState } from 'react';

function InputBar(props){

    const [url, setUrl] = useState("")
    const [customExpiry, setcustomExpiry] = useState("")

    const handleRegisterUrl = () =>{
        {   
            console.log("here1",props);
            fetch("http://localhost:3000/url-shortener/api/shorten",
                {
                    method:"POST",
                    body: JSON.stringify({
                        "url": url,
                        "expiry":parseInt(customExpiry),
                    }),
                    headers:{
                        "Content-Type":"application/json",
                        "Authorization": "Bearer " + localStorage.getItem("userToken")
                    }
                }
            ).then(response => {
                // Check if the request was successful (status code 2xx)
                if (response.ok) {
                return response.json(); // Parse the response body as JSON
                } else if (response.status === 400) {
                    // If status code is 400, handle the error
                    throw new Error(`Bad Request (${response.status}): The URL or custom expiry is invalid.`);
                } else {
                // Handle other error status codes if needed
                throw new Error(`Network response was not ok (${response.status})`);
                }
            })
            .then(data => {
                // Use the fetched data
                // Create a new array with the updated data
                const responseObject = {
                    "shortenedUrl": data.short,
                    "urlsAnalytics": {
                        "url":data.url,
                        "urlHits":0,
                        "ttl":data.expiry,
                    }
                }
                // Initialize the updatedResponseObjects with responseObject
                let updatedResponseObjects = [responseObject];
                console.log("props.urlDetailRows");
                console.log(props.urlDetailRows);
                console.log("props.urlDetailRows");
                // If props.urlDetailRows is not empty, append the new object to it
                if (props.urlDetailRows && props.urlDetailRows.length > 0) {
                    updatedResponseObjects = [...props.urlDetailRows, responseObject];
                }
                // Update the state using setUrlDetailRows with the new array
                props.setUrlDetailRows(updatedResponseObjects);
                setUrl("");
                setcustomExpiry("");
            })
            .catch(error => {
                // Handle any errors that occurred during the fetch
                console.error('Fetch error:', error);
                window.alert(error.message);
            })
        }
    }

    return (
        <div style={{marginTop:50, marginBottom:50}}>
            <div style={{"display":"flex",justifyContent:"center"}}>
                <Typography variant={"h6"}>
                    Test With Urls and ttls.
                </Typography>
            </div>
            <div style={{display:"flex", justifyContent:"center",marginTop:50, marginBottom:50}}>
                <TextField
                    id="text-field-1"
                    label="URL"
                    variant="outlined"
                    value={url}
                    style={{marginRight:20}}

                    onChange={(evant) => {
                        let elemt = evant.target;
                        setUrl(elemt.value);
                    }}
                />
                <TextField
                    id="text-field-2"
                    label="TTL in mins (default 24 hours)"
                    variant="outlined"
                    value={customExpiry}
                    style={{marginRight:40}}

                    onChange={(evant) => {
                        let elemt = evant.target;
                        setcustomExpiry(elemt.value);
                    }}
                />
                <Button 
                size={"large"} 
                variant="contained"
                onClick={handleRegisterUrl}
                >
                    Register URL
                </Button>
            </div>
        </div>
    )
}

export default InputBar