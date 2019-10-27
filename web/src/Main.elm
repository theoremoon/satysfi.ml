module Main exposing (main)

import Browser
import Browser.Navigation
import Html exposing (div, text)
import Url exposing (Url)



-- MAIN


main =
    Browser.application
        { init = init
        , update = update
        , view = view
        , subscriptions = subscriptions
        , onUrlRequest = OnUrlRequest
        , onUrlChange = OnUrlChange
        }



-- MODEL


type alias Model =
    {}


type Msg
    = OnUrlChange Url
    | OnUrlRequest Browser.UrlRequest



-- INIT


init : () -> Url -> Browser.Navigation.Key -> ( Model, Cmd Msg )
init _ url navKey =
    ( Model
    , Cmd.none
    )



-- UPDATE


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    ( model
    , Cmd.none
    )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- VIEW


view : Model -> Browser.Document Msg
view model =
    { title = "Hello World"
    , body = [ div [] [ text "Hello World" ] ]
    }
