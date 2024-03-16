import Button from '@mui/material/Button';
import TextField from "@mui/material/TextField";
import {Typography} from "@mui/material";
import { useState } from 'react';

function InputBar(props){

    const [url, setUrl] = useState("")
    const [customExpiry, setcustomExpiry] = useState("")

    const handleRegisterUrl = async () => {
        try {
            console.log("here1", props);
            const response = await fetch("http://localhost:3000/url-shortener/api/shorten", {
                method: "POST",
                body: JSON.stringify({
                    "url": url,
                    "expiry": parseInt(customExpiry),
                }),
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": "Bearer " + localStorage.getItem("userToken")
                }
            });
    
            // Check if the request was successful (status code 2xx)
            if (response.ok) {
                const data = await response.json();
                const responseObject = {
                    "shortenedUrl": data.short,
                    "urlsAnalytics": {
                        "url": data.url,
                        "urlHits": 0,
                        "ttl": data.expiry,
                    }
                }
                let updatedResponseObjects = [responseObject];
                console.log("props.urlDetailRows");
                console.log(props.urlDetailRows);
                console.log("props.urlDetailRows");
                if (props.urlDetailRows && props.urlDetailRows.length > 0) {
                    updatedResponseObjects = [...props.urlDetailRows, responseObject];
                }
                props.setUrlDetailRows(updatedResponseObjects);
                setUrl("");
                setcustomExpiry("");
            } else if (response.status === 400) {
                throw new Error(`Bad Request (${response.status}): The URL or custom expiry is invalid.`);
            } else {
                throw new Error(`Network response was not ok (${response.status})`);
            }
        } catch (error) {
            console.error('Fetch error:', error);
            window.alert(error.message);
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