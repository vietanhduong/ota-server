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

const getHost = () => {
    return process.env.REACT_APP_HOST || window.location.origin;
}

const client = axios.create({
    baseURL: `${getHost()}/api/v1`
});

function App() {
    const [data, setData] = React.useState([]);

    const fetchData = React.useCallback(() => {
        client.get('/profiles').then((response) => {
            setData(response.data || []);
        }).catch(e => console.log(e));
    }, []);

    const renderGitInfo = (metadata) => {
        metadata = metadata || {}
        // just ignore if metadata does not contains `repo`, `commit` or `pr_number` key
        if (!('repo' in metadata) || (!('commit' in metadata) && !('pr_number' in metadata))) return "";
        // if `pr_number` is null => render commit hash
        if (!metadata.pr_number) {
            return (<span className={'metadata-attribute'}>
                <b className={'mr5'}>commit:</b>
                <Link href={`${metadata.repo}/commit/${metadata.commit}`}
                      target='_blank'>{metadata.commit.substring(0, 6)}</Link>
            </span>)
        }
        // else render pr number
        return (<span className={'metadata-attribute'}>
            <b className={'mr5'}>pr:</b>
            <Link href={`${metadata.repo}/pull/${metadata.pr_number}`}
                  target='_blank'>#{metadata.pr_number}</Link>
        </span>)

    }

    const renderNoData = (d) => {
        if (d.length > 0) return '';
        return (<ListItem>
            <ListItemText primary={'No data available'}
                          style={{textAlign: "center", fontStyle: "italic", color: "gray"}}/>
        </ListItem>)
    };

    React.useEffect(() => {
        fetchData();
    }, [fetchData]);


    return (
        <Box style={{}}>
            <Box style={{maxWidth: 680, margin: "0 auto"}}>
                <List disablePadding component={Paper} style={{marginBottom: 10}} variant="outlined">
                    <ListItem>
                        <ListItemText primary={<div style={{display: "flex", justifyContent: "start"}}>
                            <img alt="" style={{width: 30}} src={"apple-icon-57x57.png"}/>
                            <Typography variant='h5' style={{paddingTop: 3}}>Over-The-Air Server</Typography>
                        </div>}/>
                    </ListItem>
                </List>
                <List disablePadding component={Paper} variant="outlined">
                    {renderNoData(data)}
                    {data.map((item, index) => (
                        <ListItem key={item.profile_id} divider={index + 1 !== data.length} style={{paddingRight: 120}}>
                            <ListItemText primary={`#${item.profile_id}: ${item.app_name}`}
                                          secondary={
                                              <>
                                                  <span className={'metadata-attribute'}>
                                                      <b className={'mr5'}>version:</b> {item.version}
                                                  </span>
                                                  <span className={'metadata-attribute'}>
                                                      <b className={'mr5'}>build:</b> {item.build}
                                                  </span>
                                                  {renderGitInfo(item.metadata)}
                                              </>
                                          }/>
                            <ListItemSecondaryAction>
                                <Link
                                    href={`itms-services://?action=download-manifest&amp;url=${getHost()}/api/v1/profiles/ios/${item.profile_id}/manifest.plist`}
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
