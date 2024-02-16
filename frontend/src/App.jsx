import { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Signup from './Signup.jsx';
import Signin from './Signin.jsx';
import UrlTable  from "./UrlTable";
import HomePage from './HomePage.jsx';
import BeanUrlAppBar from "./BeanAppbar.jsx";
import InputBar from './InputBar.jsx';
import axios from 'axios';

function App() {
    const [authenticated, setAuthenticated] = useState(false);
    const [urlDetailRows, setUrlDetailRows] = useState([]);

    useEffect(() => {
        // Use your authentication logic here
        axios.get("http://localhost:3002/user/me", {
            headers: {
                Authorization: "Bearer " + localStorage.getItem("token")
            }
        })
        .then(response => {
            setAuthenticated(true);
        })
        .catch(error => {
            setAuthenticated(false);
        });
    }, []);

    console.log("authenticated",authenticated);

    useEffect(()=>{
      fetch("http://localhost:3002/url-shortener/resolutions/analytics", { 
        method: "GET",
        headers:{
            "Content-Type":"application/json",
            "Authorization": "Bearer " + localStorage.getItem("token")
        }
     })
      .then(response => response.json())
      .then(data => {
        // Log the response
        console.log("API Response:", data);
    
        // Set the response data to your state or variable
        setUrlDetailRows(data);
      })
      .catch(error => {
        // Handle errors if any
        console.error("Error fetching data:", error);
      });
    
    },[])

    return (
        <div className="app-container"> {/* Wrap the entire app */}
            <Router>
                <BeanUrlAppBar/>
                <Routes>
                    <Route path={"/signup"} element={<Signup />} />

                    {/* Only allow access to the signin route if not authenticated */}
                    {!authenticated && <Route path={"/signin"} element={<Signin />} />}
                    {/* If Authenticated user tries to go to signin route manually.. route him/her to HomePage */}
                    {authenticated && <Route path={"/signin"} element={<HomePage urlDetailRows={urlDetailRows} setUrlDetailRows={setUrlDetailRows} />} />}
                    
                    {/* Secure routes */}
                    {authenticated ? (
                        <>
                        <Route path={"/work"} element={
                            <div>
                                <InputBar urlDetailRows={urlDetailRows} setUrlDetailRows={setUrlDetailRows} />
                                <UrlTable urlDetailRows={urlDetailRows} setUrlDetailRows={setUrlDetailRows} />
                            </div>
                        } />
                        </>
                    ) : null}

                    {/* Redirect to home if trying to access secure routes while not authenticated */}
                    {!authenticated && <Route path={"/*"} element={<HomePage urlDetailRows={urlDetailRows} setUrlDetailRows={setUrlDetailRows} />} />}
                    
                </Routes>
            </Router>
        </div>
    );
}

export default App;
