package logging

import (
    "io/ioutil"
    "fmt"
    "log"
    "io"
    "os"
)

var (
    log_DIR string
    file_PREFIX string
    max_LINES int
    use_COLOR bool
    output_HANDLE io.Writer // this doesn't seem to be used?
    log_FLAGS = log.Ldate|log.Ltime|log.Llongfile|log.LUTC

    lines int
    cur_file string

    info_logger    *log.Logger
    enter_logger   *log.Logger
    warning_logger *log.Logger
    error_logger   *log.Logger
)

// max_lines = 0 will disable file output, < 0 will disable all output
func Init(log_dir string, file_prefix string, max_lines int, use_color bool) {
    log_DIR = log_dir
    file_PREFIX = file_prefix
    max_LINES = max_lines
    use_COLOR = use_color

    info_logger = log.New(nil, "[INFO]: ", log_FLAGS)
    if use_COLOR { // hehe i hate this but its the best way
        enter_logger = log.New(nil, "\033[32m[ENTER]: ", log_FLAGS) // why nil?
        warning_logger = log.New(nil, "\033[33m[WARNING]: ", log_FLAGS)
        error_logger = log.New(nil, "\033[31m[ERROR]: ", log_FLAGS)
    } else {
        enter_logger = log.New(nil, "[ENTER]: ", log_FLAGS)
        warning_logger = log.New(nil, "[WARNING]: ", log_FLAGS)
        error_logger = log.New(nil, "[ERROR]: ", log_FLAGS)
    }

    cur_file = "" // so we always intend to have a file, from what I can tell
    if max_LINES > 0 {
        lines = max_LINES
        newLine() // forward declerations are great...
    }
}

func newLine() { // when I think of new line, I think of \n. Maybe \033[0m. This is also a ringbuffer operation, and file manipulation function.
    lines++
    if lines >= max_LINES {
        lines = 0
        if temp_file, err := ioutil.TempFile(log_DIR, file_PREFIX); err != nil {
            cur_file = ""
            Errorf("Could not create log file: "+err.Error())
        } else {
            defer temp_file.Close()
            cur_file = temp_file.Name()

            warning_logger.SetOutput(os.Stdout)
            if use_COLOR {
                warning_logger.Output(1, "New log file: "+cur_file+"\033[0m")
            } else {
                warning_logger.Output(1, "New log file: "+cur_file)
            }
        }
    }
}

func print_helper(logger *log.Logger, msg string) {
    if max_LINES < 0 {
        return
    }

    if cur_file != "" {
        if f, err := os.OpenFile(cur_file, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666); err != nil {
            if use_COLOR {
                error_logger.Output(1, "Could not open log file, "+cur_file+": "+err.Error()+"\033[0m")
            } else {
                error_logger.Output(1, "Could not open log file, "+cur_file+": "+err.Error())
            }
        } else {
            defer f.Close()

            logger.SetOutput(f) // one letter variable names are...
            logger.Output(3, msg)

            newLine()
        }
    }

    logger.SetOutput(os.Stdout)
    logger.Output(3, msg)
}

// so these are error levels but...

func Printf(format string, v ...interface{}) {
    msg := fmt.Sprintf(format, v ...)
    print_helper(info_logger, msg)
}

func Enterf(format string, v ...interface{}) {
    msg := fmt.Sprintf(format, v ...)
    if use_COLOR {
        msg += "\033[0m" // yas
    }

    print_helper(enter_logger, msg)
}

func Warningf(format string, v ...interface{}) {
    msg := fmt.Sprintf(format, v ...)
    if use_COLOR {
        msg += "\033[0m"
    }

    print_helper(warning_logger, msg)
}

func Errorf(format string, v ...interface{}) {
    msg := fmt.Sprintf(format, v ...)
    if use_COLOR {
        msg += "\033[0m"
    }

    print_helper(error_logger, msg)
}

func Check(err error) (b bool) { // is b by default false? what is this?
    if err != nil {
        msg := err.Error()
        if use_COLOR {
            msg += "\033[0m"
        }

        print_helper(error_logger, msg)

        b = true
    }
    return
}
