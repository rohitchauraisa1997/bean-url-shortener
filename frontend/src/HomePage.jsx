import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';

const HomePage = () => {
  return (
    <div>
      <Container maxWidth="md" sx={{ mt: 4 }}>
        <section>
          <Typography variant="h5" fontWeight="bold" gutterBottom>
            Bean URL Shortener
          </Typography>
          <Typography variant="body1">
            This is a simple URL shortener service that allows you to shorten long URLs into shorter, more manageable links. It's built using React as frontend, Bean as the backend server, and MySQL/Redis as the database.
          </Typography>
          <Typography variant="body1">
            You must have commonly observed these shortened links when you share links across social apps or posts on LinkedIn.
          </Typography>
          <Typography variant="h6" fontWeight="bold" gutterBottom>
            Features
          </Typography>
          <ol>
            <li>
              <Typography variant="body1">
                Shorten long URLs into concise, easy-to-share links.
              </Typography>
            </li>
            <li>
              <Typography variant="body1">
                Redirect users to the original URL when they access the shortened link.
              </Typography>
            </li>
            <li>
              <Typography variant="body1">
                Track the number of clicks on each shortened link.
              </Typography>
            </li>
            <li>
              <Typography variant="body1">
                Prevents DDOS and bot attacks by allocating a QUOTA for each user.
              </Typography>
            </li>
            <li>
              <Typography variant="body1">
                Admin UI to track all users analytics.
              </Typography>
            </li>
          </ol>
        </section>
      </Container>
    </div>
  );
};

export default HomePage;