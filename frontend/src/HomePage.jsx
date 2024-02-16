import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';

const HomePage = () => {
  return (
    <div>
      <Container maxWidth="md" sx={{ mt: 4 }}>
        <section>
          <Typography variant="body1" fontWeight="bold">
            Bean URL Shortener
          <br />
          </Typography>
        </section>
      </Container>
    </div>
  );
};

export default HomePage;