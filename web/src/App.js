import React from 'react';
import {
    Box,
    Button,
    Link,
    List,
    ListItem,
    ListItemSecondaryAction,
    ListItemText,
    Paper,
    Typography
} from '@material-ui/core';
import axios from 'axios';

const client = axios.create({
    baseURL: process.env.REACT_APP_API_V1,
});

function App() {
    const [data, setData] = React.useState([]);

    const fetchData = React.useCallback(() => {
        client.get('/profiles').then((response) => {
            setData(response.data);
        });
    }, []);

    const renderGitCommit = (metadata) => {
        metadata = metadata || {}
        // just ignore if metadata does not contain both `repo` and `commit` key
        if (!('repo' in metadata) || !('commit' in metadata)) return ''
        // display first 6 letter of commit hash
        return (<Box style={{display: "inline-flex"}}>
            <b style={{marginRight: 5}}>commit:</b>
            <Link href={`${metadata.repo}/commit/${metadata.commit}`}
                  target='_blank'>{metadata.commit.substring(0, 6)}</Link>
        </Box>)
    }

    React.useEffect(() => {
        fetchData();
    }, [fetchData]);


    return (
        <Box style={{}}>
            <Box style={{maxWidth: 680, margin: "0 auto"}}>
                <List disablePadding component={Paper} style={{marginBottom: 10}}>
                    <ListItem divider>
                        <ListItemText primary={<div style={{display: "flex", justifyContent: "start"}}>
                            <img alt="" style={{width: 30}} src={"apple-icon-57x57.png"}/>
                            <Typography variant='h5' style={{paddingTop: 3}}>Over-The-Air Server</Typography>
                        </div>}/>
                    </ListItem>
                </List>
                <List disablePadding component={Paper}>
                    {data.map((item, index) => (
                        <ListItem key={item.profile_id} divider={index + 1 !== data.length} style={{paddingRight: 120}}>
                            <ListItemText primary={item.app_name}
                                          secondary={
                                              <>
                                                  <b>version:</b> {item.version} {' '}
                                                  <b>build:</b> {item.build} {' '} {renderGitCommit(item.metadata)}
                                              </>
                                          }/>
                            <ListItemSecondaryAction>
                                <Link
                                    href={`itms-services://?action=download-manifest&amp;url=${process.env.REACT_APP_API_V1}/profiles/ios/${item.profile_id}/manifest.plist`}
                                    target='_blank'
                                >
                                    <Button disableElevation variant='contained' color='primary'
                                            style={{borderRadius: 18}}>
                                        <Typography variant='body2'>GET</Typography>
                                    </Button>
                                </Link>
                            </ListItemSecondaryAction>
                        </ListItem>
                    ))}
                </List>
            </Box>
        </Box>
    );
}

export default App;
