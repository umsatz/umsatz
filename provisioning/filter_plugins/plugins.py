from ansible import errors

def split_string(string, seperator=' '):
    return string.split(seperator)

def seq_take(seq, number):
    return seq[0:number]

class FilterModule(object):
    ''' A filter to split a string into a list. '''
    def filters(self):
        return {
            'split': split_string,
            'take': seq_take
        }