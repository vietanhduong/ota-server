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
        // just ignore if metadata does not contain both `repo` and `commit` key
        if (!('repo' in metadata) || !('commit' in metadata)) return ''
        // display first 6 letter of commit hash
        return (<Link href={`${metadata.repo}/commit/${metadata.commit}`}
                      target='_blank'>(#{metadata.commit.substring(0, 6)})</Link>)
    }

    React.useEffect(() => {
        fetchData();
    }, [fetchData]);


    return (
        <Box>
            <List disablePadding component={Paper}>
                <ListItem divider>
                    <ListItemText primary={<Typography variant='h4'><img alt="" style={{width: 30}}
                                                                         src={"apple-icon-57x57.png"}/> Over-The-Air
                        Server</Typography>}/>
                </ListItem>
                {data.map((item, index) => (
                    <ListItem key={item.profile_id} divider={index + 1 !== data.length} style={{paddingRight: 120}}>
                        <ListItemText primary={item.app_name}
                                      secondary={
                                          <>
                                              <b>Version:</b> {item.version} {' '}
                                              <b>Build:</b> {item.build} {' '} {renderGitCommit(item.metadata)}
                                          </>
                                      }/>
                        <ListItemSecondaryAction>
                            <Link
                                href={`itms-services://?action=download-manifest&amp;url=${process.env.REACT_APP_API_V1}/profiles/ios/${item.profile_id}/manifest.plist`}
                                target='_blank'
                            >
                                <Button disableElevation variant='contained' color='primary' style={{borderRadius: 18}}>
                                    <Typography variant='body2'>INSTALL</Typography>
                                </Button>
                            </Link>
                        </ListItemSecondaryAction>
                    </ListItem>
                ))}
            </List>
        </Box>
    );
}

export default App;
