import React, { useState, useEffect } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { cloneDeep, get, concat, isUndefined } from "lodash";
import Select from '@material-ui/core/Select';
import MenuItem from '@material-ui/core/MenuItem';
import Paper from '@material-ui/core/Paper';
import Tooltip from '@material-ui/core/Tooltip';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TablePagination from '@material-ui/core/TablePagination';
import TableRow from '@material-ui/core/TableRow';
import FormControl from '@material-ui/core/FormControl';
import InputLabel from '@material-ui/core/InputLabel';
import Fab from '@material-ui/core/Fab';
import IconButton from '@material-ui/core/IconButton';
import Backdrop from '@material-ui/core/Backdrop';
import Dialog from '@material-ui/core/Dialog';
import Fade from '@material-ui/core/Fade';
import DoneIcon from '@material-ui/icons/Done';
import UndoIcon from '@material-ui/icons/Undo';
import AddIcon from '@material-ui/icons/Add';
import DeleteIcon from '@material-ui/icons/Delete';
import { TaskAPI, Task } from "./../../services/tasks";
import { ListAPI } from "./../../services/list";
import { ProjectAPI } from "./../../services/project";
import { BacklogViewAPI, BacklogList, BacklogTask, BacklogProject } from "../../services/views/backlog";
import { TaskModal } from './new-task-modal/modal';

const useStyles = makeStyles((theme) => ({
  root: {
    width: '100%',
  },
  container: {
    height: '100%',
  },
  formControl: {
    margin: theme.spacing(1),
    minWidth: 120,
  },
  modal: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
  paper: {
    backgroundColor: theme.palette.background.paper,
    border: '2px solid #000',
    boxShadow: theme.shadows[5],
    padding: theme.spacing(2, 4, 3),
  },
}));


interface IProps {
  TaskAPI: TaskAPI;
  ListAPI: ListAPI;
  ProjectAPI: ProjectAPI;
  BacklogViewAPI: BacklogViewAPI;
}

interface Row extends BacklogTask { }
interface Column {
  title: string;
  minWidth: number;
}

interface IState {
  lists: BacklogList[];
  projects: BacklogProject[];
  columns: Column[];
}

function makeColumn(title: string): Column {
  return {
    title,
    minWidth: 30
  }
}
interface SelectionChanged {
  id: string;
  name: string;
}

interface Selection {
  id: string;
  name: string;
}
function makeListSelection(name: string, selected: Selection, row: Row, items: Selection[], className: string, onChange: (source: SelectionChanged, destination: SelectionChanged) => void) {
  const menuItems = concat([{ name: 'Empty', id: "-1" }], items).map(({ name, id }) => {
    return (
      <MenuItem key={id} value={id}>{name}</MenuItem>
    )
  })
  return (
    <FormControl className={className}>
      <InputLabel >{name}</InputLabel>
      <Select
        value={selected.id || "-1"}
        onChange={(event: React.ChangeEvent<{ name?: string; value: unknown }>, child: React.ReactNode) => {
          const id = get(child, 'props.value', "-1");
          const name = get(child, 'props.children', "Empty");
          const destination = {
            id: id === "-1" ? "" : id,
            name: name === "Empty" ? "" : name,
          }
          const source = {
            id: selected.id,
            name: selected.name,
          };
          onChange(source, destination)
        }}
      >
        {menuItems}
      </Select>
    </FormControl>
  )
}

function makeTaskCompletionButton(row: Row, onClick: (action: string) => void) {
  const icon = row.status.completed ? <UndoIcon></UndoIcon> : <DoneIcon></DoneIcon>
  const action = !row.status.completed ? 'Complete' : 'Reopen'
  return (
    <Tooltip title={action} aria-label={action}>
      <IconButton aria-label="toggelCompletion" color="primary" id={row.metadata.id} onClick={() => {
        onClick(action)
      }}>
        {icon}
      </IconButton>
    </Tooltip>
  )
}

export function BacklogPage(props: IProps) {
  const classes = useStyles();
  const [taskModal, setTaskModelOpen] = useState(false);
  const handleTaskModalOpen = () => {
    setTaskModelOpen(true);
  };

  const handleTaskModalClose = (data: any) => {
    setTaskModelOpen(false);
  };
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [rows, updateRows] = useState<Row[]>([]);
  const handleChangePage = (event: unknown, newPage: number) => {
    setPage(newPage);
  };
  const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(+event.target.value);
    setPage(0);
  };
  const [state, setState] = useState<IState>(() => {
    return {
      lists: [],
      projects: [],
      columns: [],
    }
  });

  useEffect(() => {
    (async () => {
      const view = await props.BacklogViewAPI.get()
      const state: IState = {
        lists: view.lists,
        projects: view.projects,
        columns: [
          makeColumn('Actions'),
          makeColumn('Title'),
          makeColumn('Description'),
          makeColumn('List'),
          makeColumn('Project'),
        ],
      }
      setState(state)
      updateRows(view.tasks)
    })()
  }, [props.ListAPI, props.TaskAPI, props.BacklogViewAPI])

  const onUpdate = async (newData: Row, oldData?: Row): Promise<any> => {
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve();
        updateRows((prevState) => {
          let index;
          for (let i = 0; i < prevState.length; i++) {
            const element = prevState[i];
            if (element.metadata.id === newData.metadata.id) {
              index = i
            }
          }
          // index does not found, return previous state
          if (isUndefined(index)) {
            return prevState
          };
          const data = [...prevState];
          data[index] = newData;
          return data;
        });
      }, 600);
    });
  }

  const addTask = async (row: Row): Promise<void> => {
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve();
        updateRows((prevState) => {
          const data = [...prevState];
          data.splice(0, 0, row)
          return data;
        });
      }, 600);
    });
  }

  const deleteTask = async (row: Row): Promise<void> => {
    await props.TaskAPI.remove(row.metadata.id)
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve();
        updateRows((prevState) => {
          let index;
          for (let i = 0; i < prevState.length; i++) {
            const element = prevState[i];
            if (element.metadata.id === row.metadata.id) {
              index = i
            }
          }
          // index does not found, return previous state
          if (isUndefined(index)) {
            return prevState
          };
          const data = [...prevState];
          data.splice(index, 1)
          return data;
        });
      }, 600);
    });
  }

  return (
    <Paper className={classes.root}>
      <TableContainer className={classes.container}>
        <Table stickyHeader aria-label="sticky table" size={'small'}>
          <TableHead>
            <TableRow>
              {state.columns.map((column, index) => (
                <TableCell
                  key={index}

                  style={{ minWidth: column.minWidth }}
                >
                  {column.title}
                </TableCell>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {rows.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage).map((row, index) => {
              return (
                <TableRow hover role="checkbox" tabIndex={-1} key={index}>
                  <TableCell key={'c-0'}>
                    {makeTaskCompletionButton(row, async (action: string) => {
                      if (action === 'Reopen') {
                        await props.TaskAPI.reopen(row.metadata.id);
                      }

                      if (action === 'Complete') {
                        await props.TaskAPI.complete(row.metadata.id);
                      }
                      const newRow = cloneDeep(row)
                      newRow.status.completed = !newRow.status.completed
                      onUpdate(newRow, row)
                    })}
                    <Tooltip title={"Delete"} aria-label={"delete"}>
                      <IconButton aria-label="toggelCompletion" color="primary" id={row.metadata.id} onClick={() => {
                        deleteTask(row)
                      }}>
                        <DeleteIcon />
                      </IconButton>
                    </Tooltip>
                  </TableCell>
                  <TableCell key={'c-1'}>
                    {row.metadata.name}
                  </TableCell>
                  <TableCell key={'c-2'}>
                    {row.metadata.description}
                  </TableCell>
                  <TableCell key={'c-3'}>
                    {makeListSelection('List', { id: row.list.id || '', name: row.list.name || '' }, row, state.lists, classes.formControl, async (source: SelectionChanged, destination: SelectionChanged) => {
                      await props.ListAPI.moveTasks(source.id, destination.id, [row.metadata.id])
                      const newRow = cloneDeep(row)
                      newRow.list.id = destination.id
                      newRow.list.name = destination.name
                      onUpdate(newRow, row)
                    })}
                  </TableCell>
                  <TableCell key={'c-4'}>
                    {makeListSelection('Project', { id: row.project.id || '', name: row.project.name || '' }, row, state.projects, classes.formControl, async (source: SelectionChanged, destination: SelectionChanged) => {
                      await props.ProjectAPI.addTasks(destination.id, [row.metadata.id])
                      const newRow = cloneDeep(row)
                      newRow.project.id = destination.id
                      newRow.project.name = destination.name
                      onUpdate(newRow, row)
                    })}
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </TableContainer>
      <TablePagination
        rowsPerPageOptions={[10, 25, 100]}
        component="div"
        count={1}
        rowsPerPage={rowsPerPage}
        page={page}
        onChangePage={handleChangePage}
        onChangeRowsPerPage={handleChangeRowsPerPage}
      />
      <Fab
        color="primary"
        style={{ position: 'fixed', right: '15px', bottom: '15px' }}
        onClick={handleTaskModalOpen}
      >
        <AddIcon />
      </Fab>
      <Dialog
        aria-labelledby="transition-modal-title"
        aria-describedby="transition-modal-description"
        className={classes.modal}
        open={taskModal}
        onClose={handleTaskModalClose}
        closeAfterTransition
        BackdropComponent={Backdrop}
        BackdropProps={{
          timeout: 500,
        }}
      >
        <Fade in={taskModal}>
          <TaskModal TaskAPI={props.TaskAPI} Callback={(task: Task) => {
            setTaskModelOpen(false)
            const row = {
              ...task,
              list: {
                id: '',
                name: '',
              },
              project: {
                id: '',
                name: '',
              }
            }
            addTask(row)
          }}></TaskModal>
        </Fade>
      </Dialog>
    </Paper>
  );
}
